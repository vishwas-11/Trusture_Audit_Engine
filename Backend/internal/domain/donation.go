package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type DonationStatus string

const (
	DonationPending DonationStatus = "PENDING_AUDIT"
	DonationValid   DonationStatus = "VALIDATED"
	DonationFlagged DonationStatus = "FLAGGED"
)

type Donation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TxHash    string             `bson:"tx_hash" json:"tx_hash"`
	Donor     string             `bson:"donor_wallet" json:"donor_wallet"`
	NGO       string             `bson:"ngo_wallet" json:"ngo_wallet"`
	Amount    float64            `bson:"amount" json:"amount"`
	Purpose   string             `bson:"purpose" json:"purpose"`
	Status    DonationStatus     `bson:"status" json:"status"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
}
