package service

import (
	"time"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/repository"
)

type DonationService struct {
	donationRepo repository.DonationRepository
}

func NewDonationService(repo repository.DonationRepository) *DonationService {
	return &DonationService{
		donationRepo: repo,
	}
}

// This will later be triggered by blockchain events
func (s *DonationService) RecordDonation(
	txHash string,
	donorWallet string,
	ngoWallet string,
	amount float64,
	purpose string,
) (*domain.Donation, error) {

	donation := &domain.Donation{
		TxHash:    txHash,
		Donor:     donorWallet,
		NGO:       ngoWallet,
		Amount:    amount,
		Purpose:   purpose,
		Status:    domain.DonationPending,
		CreatedAt: time.Now().Unix(),
	}

	err := s.donationRepo.Create(donation)
	if err != nil {
		return nil, err
	}

	return donation, nil
}

func (s *DonationService) GetDonationsByNGOWallet(ngoWallet string) ([]domain.Donation, error) {
	return s.donationRepo.FindByNGOWallet(ngoWallet)
}

func (s *DonationService) ListAllRecent(limit int64) ([]domain.Donation, error) {
	return s.donationRepo.FindAllRecent(limit)
}

func (s *DonationService) GetByID(id string) (*domain.Donation, error) {
	return s.donationRepo.FindByID(id)
}

func (s *DonationService) UpdateStatus(id string, status domain.DonationStatus) error {
	return s.donationRepo.UpdateStatus(id, status)
}
