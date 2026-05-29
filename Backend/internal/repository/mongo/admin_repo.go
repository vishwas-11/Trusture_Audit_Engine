package mongo

import (
	"context"
	"errors"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	collection *mongo.Collection
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{collection: db.Collection("admins")}
}

func (r *AdminRepository) Create(admin *domain.AdminUser) error {
	res, err := r.collection.InsertOne(context.Background(), admin)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		admin.ID = oid
	}
	return nil
}

func (r *AdminRepository) FindByEmail(email string) (*domain.AdminUser, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	var u domain.AdminUser
	err := r.collection.FindOne(context.Background(), bson.M{"email": normalized}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
