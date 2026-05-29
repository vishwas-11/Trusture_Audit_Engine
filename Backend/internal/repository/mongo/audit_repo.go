package mongo

import (
	"context"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuditRepository struct {
	collection *mongo.Collection
}

func NewAuditRepository(db *mongo.Database) *AuditRepository {
	return &AuditRepository{
		collection: db.Collection("audit_results"),
	}
}

func (r *AuditRepository) Create(result *domain.AuditResult) error {
	res, err := r.collection.InsertOne(context.Background(), result)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		result.ID = oid
	}
	return nil
}

func (r *AuditRepository) FindLatestByDonationID(donationID string) (*domain.AuditResult, error) {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(donationID))
	if err != nil {
		return nil, err
	}

	var result domain.AuditResult
	err = r.collection.FindOne(
		context.Background(),
		bson.M{"donation_id": objectID},
		options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AuditRepository) FindByDonationID(donationID string) ([]domain.AuditResult, error) {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(donationID))
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"donation_id": objectID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	results := make([]domain.AuditResult, 0)
	for cursor.Next(context.Background()) {
		var result domain.AuditResult
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
