package service

import (
	"time"

	"github.com/vishwas-11/trusture-backend/internal/domain"
)

type RiskService struct {
	geo    *VendorService
	vendor *VendorService
}

func NewRiskService() *RiskService {
	return &RiskService{
		geo:    NewVendorService(),
		vendor: NewVendorService(),
	}
}

func (r *RiskService) AssessRisk(
	donationID domain.Donation,
	vendorGST string,
	geoMismatch bool,
) *domain.RiskAssessment {

	score := 0.0
	reasons := []string{}

	if geoMismatch {
		score += 0.4
		reasons = append(reasons, "GEO_LOCATION_MISMATCH")
	}

	if !r.vendor.VerifyGST(vendorGST) {
		score += 0.5
		reasons = append(reasons, "UNVERIFIED_VENDOR")
	}

	level := domain.RiskLow
	if score > 0.7 {
		level = domain.RiskHigh
	} else if score > 0.4 {
		level = domain.RiskMedium
	}

	return &domain.RiskAssessment{
		DonationID: donationID.ID,
		Level:      level,
		Score:      score,
		Reasons:    reasons,
		CreatedAt:  time.Now().Unix(),
	}
}
