package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole string

const (
	RoleDonor UserRole = "DONOR"
	RoleNGO   UserRole = "NGO"
	RoleAdmin UserRole = "ADMIN"
)

type KYCStatus string

const (
	KYCPending  KYCStatus = "PENDING"
	KYCVerified KYCStatus = "VERIFIED"
	KYCRejected KYCStatus = "REJECTED"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Role      UserRole           `bson:"role"`
	Wallet    string             `bson:"wallet"`
	DID       string             `bson:"did"`
	KYCStatus KYCStatus          `bson:"kyc_status"`
	Latitude  float64            `bson:"latitude,omitempty"`
	Longitude float64            `bson:"longitude,omitempty"`
	CreatedAt int64              `bson:"created_at"`
}
