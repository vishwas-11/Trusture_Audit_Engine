package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuditStatus string

const (
	AuditPass AuditStatus = "PASS"
	AuditFail AuditStatus = "FAIL"
)

type AuditResult struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DonationID            primitive.ObjectID `bson:"donation_id" json:"donation_id"`
	Status                AuditStatus        `bson:"status" json:"status"`
	Score                 float64            `bson:"score" json:"score"`
	Flags                 []string           `bson:"flags" json:"flags"`
	DonorNGODistanceKM    float64            `bson:"donor_ngo_distance_km,omitempty" json:"donor_ngo_distance_km,omitempty"`
	NGOProofDistanceKM    float64            `bson:"ngo_proof_distance_km,omitempty" json:"ngo_proof_distance_km,omitempty"`
	GSTVerificationSource string             `bson:"gst_verification_source,omitempty" json:"gst_verification_source,omitempty"`
	CreatedAt             int64              `bson:"created_at" json:"created_at"`
}
