package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// AdminUser is stored in MongoDB for email/password admin dashboard access.
type AdminUser struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	CreatedAt    int64              `bson:"created_at" json:"created_at"`
}
