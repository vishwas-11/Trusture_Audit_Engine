package repository

import "github.com/vishwas-11/trusture-backend/internal/domain"

type DonationRepository interface {
	Create(donation *domain.Donation) error
	FindByTxHash(txHash string) (*domain.Donation, error)
	FindByID(id string) (*domain.Donation, error)
	FindByNGOWallet(ngoWallet string) ([]domain.Donation, error)
	FindAllRecent(limit int64) ([]domain.Donation, error)
	UpdateStatus(id string, status domain.DonationStatus) error
}

type ProofRepository interface {
	Create(proof *domain.Proof) error
	FindByDonationID(donationID string) ([]domain.Proof, error)
	UpdateExtraction(proofID string, extractedText string, extractionNote string) error
}

type AdminRepository interface {
	Create(admin *domain.AdminUser) error
	FindByEmail(email string) (*domain.AdminUser, error)
}

type AuditRepository interface {
	Create(result *domain.AuditResult) error
	FindLatestByDonationID(donationID string) (*domain.AuditResult, error)
	FindByDonationID(donationID string) ([]domain.AuditResult, error)
}
