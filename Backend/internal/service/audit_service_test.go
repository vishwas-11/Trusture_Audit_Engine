package service

import (
	"context"
	"testing"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRunAuditFailsOnMissingProofs(t *testing.T) {
	svc := NewAuditService(NewVendorService(), 0.70, 100, 600, true)
	donation := &domain.Donation{ID: primitive.NewObjectID()}

	result := svc.RunAudit(context.Background(), donation, nil, AuditInput{})
	if result.Status != domain.AuditFail {
		t.Fatalf("expected FAIL, got %s", result.Status)
	}
	if result.Score >= 0.70 {
		t.Fatalf("expected score below threshold, got %f", result.Score)
	}
}

func TestRunAuditFlagsMediaWithoutGeotag(t *testing.T) {
	svc := NewAuditService(NewVendorService(), 0.70, 100, 600, true)
	donation := &domain.Donation{ID: primitive.NewObjectID()}
	proofs := []domain.Proof{
		{IPFSHash: "QmHash1", Type: domain.ProofMedia},
	}

	result := svc.RunAudit(context.Background(), donation, proofs, AuditInput{})
	if !hasFlag(result.Flags, "MEDIA_WITHOUT_GEOTAG") {
		t.Fatalf("expected MEDIA_WITHOUT_GEOTAG flag, got %v", result.Flags)
	}
}

func TestRunAuditUsesGSTFallbackFlag(t *testing.T) {
	svc := NewAuditService(NewVendorService(), 0.70, 100, 600, true)
	donation := &domain.Donation{ID: primitive.NewObjectID()}
	proofs := []domain.Proof{
		{IPFSHash: "QmHash2", Type: domain.ProofInvoice, GSTIN: "29ABCDE1234F1Z5"},
	}

	result := svc.RunAudit(context.Background(), donation, proofs, AuditInput{})
	if !hasFlag(result.Flags, "GST_UNVERIFIED_API_DOWN") {
		t.Fatalf("expected GST fallback flag, got %v", result.Flags)
	}
	if result.GSTVerificationSource != "FORMAT_FALLBACK" {
		t.Fatalf("expected FORMAT_FALLBACK source, got %s", result.GSTVerificationSource)
	}
}

func TestRunAuditDistancePenalty(t *testing.T) {
	svc := NewAuditService(NewVendorService(), 0.70, 100, 600, true)
	donation := &domain.Donation{ID: primitive.NewObjectID()}
	proofs := []domain.Proof{
		{IPFSHash: "QmHash3", Type: domain.ProofMedia, GeoTag: "tag", Latitude: 12.9716, Longitude: 77.5946},
	}

	input := AuditInput{
		DonorProfileLat: 28.6139, DonorProfileLng: 77.2090,
		NGOProfileLat: 13.0827, NGOProfileLng: 80.2707,
		ProofLat: 12.9716, ProofLng: 77.5946,
	}

	result := svc.RunAudit(context.Background(), donation, proofs, input)
	if !hasFlag(result.Flags, "DONOR_NGO_TOO_FAR_SEVERE") {
		t.Fatalf("expected severe donor-ngo distance flag, got %v", result.Flags)
	}
}

func hasFlag(flags []string, expected string) bool {
	for _, flag := range flags {
		if flag == expected {
			return true
		}
	}
	return false
}
