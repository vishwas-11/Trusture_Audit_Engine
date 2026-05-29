package service

import (
	"context"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/repository"
)

type AuditPipelineService struct {
	donationService *DonationService
	proofService    *ProofService
	auditService    *AuditService
	auditRepo       repository.AuditRepository
}

func NewAuditPipelineService(
	donationService *DonationService,
	proofService *ProofService,
	auditService *AuditService,
	auditRepo repository.AuditRepository,
) *AuditPipelineService {
	return &AuditPipelineService{
		donationService: donationService,
		proofService:    proofService,
		auditService:    auditService,
		auditRepo:       auditRepo,
	}
}

func (s *AuditPipelineService) RunForDonation(ctx context.Context, donationID string, input AuditInput) (*domain.AuditResult, error) {
	donation, err := s.donationService.GetByID(donationID)
	if err != nil {
		return nil, err
	}

	proofs, err := s.proofService.ListByDonationID(donationID)
	if err != nil {
		return nil, err
	}

	result := s.auditService.RunAudit(ctx, donation, proofs, input)
	if err := s.auditRepo.Create(result); err != nil {
		return nil, err
	}

	nextStatus := domain.DonationValid
	if result.Status == domain.AuditFail {
		nextStatus = domain.DonationFlagged
	}
	if err := s.donationService.UpdateStatus(donationID, nextStatus); err != nil {
		return nil, err
	}

	return result, nil
}
