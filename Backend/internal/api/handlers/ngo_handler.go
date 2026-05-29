package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/service"
	"github.com/vishwas-11/trusture-backend/internal/utils"
)

var ngoDonationService *service.DonationService

func InitNGOHandler(donationService *service.DonationService) {
	ngoDonationService = donationService
}

func NGODashboard(w http.ResponseWriter, r *http.Request) {
	if ngoDonationService == nil {
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"SERVICE_NOT_READY",
			"Donation service not initialized",
		)
		return
	}

	wallet := strings.TrimSpace(r.Header.Get("X-WALLET"))
	if wallet == "" {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_WALLET",
			"Wallet header is required",
		)
		return
	}

	donations, err := ngoDonationService.GetDonationsByNGOWallet(wallet)
	if err != nil {
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"FETCH_FAILED",
			"Failed to fetch NGO transactions",
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(donations); err != nil {
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"ENCODE_FAILED",
			"Failed to encode NGO transactions",
		)
		return
	}
}
