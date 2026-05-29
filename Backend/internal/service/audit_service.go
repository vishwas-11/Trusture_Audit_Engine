package service

import (
	"context"
	"strings"
	"time"

	"github.com/vishwas-11/trusture-backend/internal/domain"
)

type AuditService struct {
	vendorService *VendorService
	passThreshold float64
	distanceWarnKM float64
	distanceFailKM float64
	gstEnabled    bool
}

type AuditInput struct {
	DonorProfileLat float64
	DonorProfileLng float64
	NGOProfileLat   float64
	NGOProfileLng   float64
	ProofLat        float64
	ProofLng        float64
	RequestLat      float64
	RequestLng      float64
}

func NewAuditService(
	vendorService *VendorService,
	passThreshold float64,
	distanceWarnKM float64,
	distanceFailKM float64,
	gstEnabled bool,
) *AuditService {
	return &AuditService{
		vendorService: vendorService,
		passThreshold: passThreshold,
		distanceWarnKM: distanceWarnKM,
		distanceFailKM: distanceFailKM,
		gstEnabled:    gstEnabled,
	}
}

func (a *AuditService) RunAudit(
	ctx context.Context,
	donation *domain.Donation,
	proofs []domain.Proof,
	input AuditInput,
) *domain.AuditResult {

	var flags []string
	score := 1.0
	gstSource := "NOT_REQUIRED"

	ngoPoint, hasNGOProfile := buildGeoPoint(input.NGOProfileLat, input.NGOProfileLng)
	donorPoint, hasDonorProfile := buildGeoPoint(input.DonorProfileLat, input.DonorProfileLng)
	proofPoint, hasProofPoint := buildGeoPoint(input.ProofLat, input.ProofLng)
	requestPoint, hasRequestPoint := buildGeoPoint(input.RequestLat, input.RequestLng)

	var donorNGODistance float64
	var ngoProofDistance float64
	if hasDonorProfile && hasNGOProfile {
		donorNGODistance = DistanceKM(donorPoint, ngoPoint)
		score, flags = a.applyDistancePenalty(score, flags, donorNGODistance, "DONOR_NGO_TOO_FAR")
	} else {
		flags = append(flags, "MISSING_PROFILE_LOCATION")
		score -= 0.1
	}

	ngoReference := ngoPoint
	hasNGOReference := hasNGOProfile
	if !hasNGOReference && hasRequestPoint {
		ngoReference = requestPoint
		hasNGOReference = true
		flags = append(flags, "USING_REQUEST_LOCATION_FALLBACK")
	}

	if hasNGOReference && hasProofPoint {
		ngoProofDistance = DistanceKM(ngoReference, proofPoint)
		score, flags = a.applyDistancePenalty(score, flags, ngoProofDistance, "NGO_PROOF_TOO_FAR")
	}

	// Rule 1: Proof existence
	if len(proofs) == 0 {
		flags = append(flags, "NO_PROOFS_UPLOADED")
		score -= 0.5
	}

	ipfsSeen := make(map[string]bool)

	for _, proof := range proofs {

		// Rule 2: IPFS hash must exist
		if proof.IPFSHash == "" {
			flags = append(flags, "EMPTY_IPFS_HASH")
			score -= 0.2
		}

		// Rule 3: Duplicate proof detection
		if ipfsSeen[proof.IPFSHash] {
			flags = append(flags, "DUPLICATE_PROOF_DETECTED")
			score -= 0.3
		}
		ipfsSeen[proof.IPFSHash] = true

		// Rule 4: Geo-tag mandatory for media
		if proof.Type == domain.ProofMedia && proof.GeoTag == "" {
			flags = append(flags, "MEDIA_WITHOUT_GEOTAG")
			score -= 0.2
		}

		requiresGST := proof.Type == domain.ProofInvoice || proof.Type == domain.ProofReceipt
		if requiresGST {
			if strings.TrimSpace(proof.GSTIN) == "" {
				flags = append(flags, "MISSING_GSTIN")
				score -= 0.35
				continue
			}
			if a.gstEnabled {
				isValid, source := a.vendorService.VerifyGSTWithFallback(ctx, strings.ToUpper(strings.TrimSpace(proof.GSTIN)))
				gstSource = source
				if !isValid {
					flags = append(flags, "INVALID_GSTIN")
					score -= 0.4
				} else if source == "FORMAT_FALLBACK" {
					flags = append(flags, "GST_UNVERIFIED_API_DOWN")
					score -= 0.05
				}
			} else {
				gstSource = "DISABLED"
			}
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	status := domain.AuditPass
	if score < a.passThreshold {
		status = domain.AuditFail
	}

	return &domain.AuditResult{
		DonationID:            donation.ID,
		Status:                status,
		Score:                 score,
		Flags:                 flags,
		DonorNGODistanceKM:    donorNGODistance,
		NGOProofDistanceKM:    ngoProofDistance,
		GSTVerificationSource: gstSource,
		CreatedAt:             time.Now().Unix(),
	}
}

func buildGeoPoint(lat, lng float64) (GeoPoint, bool) {
	if lat == 0 && lng == 0 {
		return GeoPoint{}, false
	}
	return GeoPoint{Lat: lat, Lng: lng}, true
}

func (a *AuditService) applyDistancePenalty(score float64, flags []string, distanceKM float64, flagPrefix string) (float64, []string) {
	if distanceKM > a.distanceFailKM {
		flags = append(flags, flagPrefix+"_SEVERE")
		return score - 0.35, flags
	}
	if distanceKM > a.distanceWarnKM {
		flags = append(flags, flagPrefix+"_WARN")
		return score - 0.15, flags
	}
	return score, flags
}
