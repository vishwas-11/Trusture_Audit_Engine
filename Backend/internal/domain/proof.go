package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProofType string

const (
	ProofInvoice ProofType = "INVOICE"
	ProofReceipt ProofType = "RECEIPT"
	ProofMedia   ProofType = "MEDIA"
)

type Proof struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DonationID  primitive.ObjectID `bson:"donation_id" json:"donation_id"`
	IPFSHash    string             `bson:"ipfs_hash" json:"ipfs_hash"`
	Type        ProofType          `bson:"type" json:"type"`
	GeoTag      string             `bson:"geo_tag,omitempty" json:"geo_tag,omitempty"`
	Latitude    float64            `bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude   float64            `bson:"longitude,omitempty" json:"longitude,omitempty"`
	GSTIN       string             `bson:"gstin,omitempty" json:"gstin,omitempty"`
	VendorName  string             `bson:"vendor_name,omitempty" json:"vendor_name,omitempty"`
	// ExtractedText is plain text from the PDF text layer or OCR/vision (when API keys are configured).
	ExtractedText string `bson:"extracted_text,omitempty" json:"extracted_text,omitempty"`
	// ExtractionNote explains why text is empty or that nothing relevant was found.
	ExtractionNote string `bson:"extraction_note,omitempty" json:"extraction_note,omitempty"`
	UploadedAt  int64              `bson:"uploaded_at" json:"uploaded_at"`
}
