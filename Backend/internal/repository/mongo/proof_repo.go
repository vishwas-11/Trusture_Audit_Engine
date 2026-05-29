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

type ProofRepository struct {
	collection *mongo.Collection
}

func NewProofRepository(db *mongo.Database) *ProofRepository {
	return &ProofRepository{
		collection: db.Collection("proofs"),
	}
}

func (r *ProofRepository) Create(proof *domain.Proof) error {
	res, err := r.collection.InsertOne(context.Background(), proof)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		proof.ID = oid
	}
	return nil
}

func (r *ProofRepository) UpdateExtraction(proofID string, extractedText string, extractionNote string) error {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(proofID))
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"extracted_text":   extractedText,
			"extraction_note": extractionNote,
		}},
	)
	return err
}

func (r *ProofRepository) FindByDonationID(donationID string) ([]domain.Proof, error) {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(donationID))
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"donation_id": objectID},
		options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	proofs := make([]domain.Proof, 0)
	for cursor.Next(context.Background()) {
		var proof domain.Proof
		if err := cursor.Decode(&proof); err != nil {
			return nil, err
		}
		proofs = append(proofs, proof)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return proofs, nil
}
