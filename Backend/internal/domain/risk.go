package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RiskLevel string

const (
	RiskLow    RiskLevel = "LOW"
	RiskMedium RiskLevel = "MEDIUM"
	RiskHigh   RiskLevel = "HIGH"
)

type RiskAssessment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DonationID primitive.ObjectID `bson:"donation_id"`
	Level       RiskLevel          `bson:"level"`
	Score       float64            `bson:"score"`
	Reasons     []string           `bson:"reasons"`
	CreatedAt   int64              `bson:"created_at"`
}
