package service

import (
	"time"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProofService struct {
	proofRepo repository.ProofRepository
}

func NewProofService(repo repository.ProofRepository) *ProofService {
	return &ProofService{proofRepo: repo}
}

func (s *ProofService) CreateProof(
	donationID string,
	ipfsHash string,
	proofType domain.ProofType,
	geoTag string,
	latitude float64,
	longitude float64,
	gstin string,
	vendorName string,
	extractedText string,
	extractionNote string,
) (*domain.Proof, error) {
	donationObjectID, err := primitive.ObjectIDFromHex(donationID)
	if err != nil {
		return nil, err
	}

	proof := &domain.Proof{
		DonationID:     donationObjectID,
		IPFSHash:       ipfsHash,
		Type:           proofType,
		GeoTag:         geoTag,
		Latitude:       latitude,
		Longitude:      longitude,
		GSTIN:          gstin,
		VendorName:     vendorName,
		ExtractedText:  extractedText,
		ExtractionNote: extractionNote,
		UploadedAt: time.Now().Unix(),
	}

	if err := s.proofRepo.Create(proof); err != nil {
		return nil, err
	}

	return proof, nil
}

func (s *ProofService) ListByDonationID(donationID string) ([]domain.Proof, error) {
	return s.proofRepo.FindByDonationID(donationID)
}
