package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type NGO struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Registration  string             `bson:"registration"`
	FCRACompliant bool               `bson:"fcra_compliant"`
	Wallet        string             `bson:"wallet"`
	CreatedAt     int64              `bson:"created_at"`
}
