package mongo

import (
	"context"
	"errors"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DonationRepository struct {
	collection *mongo.Collection
}

func NewDonationRepository(db *mongo.Database) *DonationRepository {
	return &DonationRepository{
		collection: db.Collection("donations"),
	}
}

func (r *DonationRepository) Create(d *domain.Donation) error {
	_, err := r.collection.InsertOne(context.Background(), d)
	return err
}

func (r *DonationRepository) FindByTxHash(txHash string) (*domain.Donation, error) {
	var donation domain.Donation
	err := r.collection.FindOne(
		context.Background(),
		map[string]string{"tx_hash": txHash},
	).Decode(&donation)

	return &donation, err
}

func (r *DonationRepository) FindByID(id string) (*domain.Donation, error) {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}

	var donation domain.Donation
	err = r.collection.FindOne(
		context.Background(),
		bson.M{"_id": objectID},
	).Decode(&donation)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	return &donation, err
}

func (r *DonationRepository) FindByNGOWallet(ngoWallet string) ([]domain.Donation, error) {
	normalizedWallet := strings.ToLower(strings.TrimSpace(ngoWallet))

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{
			"$expr": bson.M{
				"$eq": bson.A{
					bson.M{"$toLower": "$ngo_wallet"},
					normalizedWallet,
				},
			},
		},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	donations := make([]domain.Donation, 0)
	for cursor.Next(context.Background()) {
		var donation domain.Donation
		if err := cursor.Decode(&donation); err != nil {
			return nil, err
		}
		donations = append(donations, donation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return donations, nil
}

func (r *DonationRepository) FindAllRecent(limit int64) ([]domain.Donation, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	out := make([]domain.Donation, 0)
	for cursor.Next(context.Background()) {
		var d domain.Donation
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *DonationRepository) UpdateStatus(id string, status domain.DonationStatus) error {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(id))
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}
