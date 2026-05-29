package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/vishwas-11/trusture-backend/internal/api/middleware"
	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/repository"
	"github.com/vishwas-11/trusture-backend/internal/service"
	"github.com/vishwas-11/trusture-backend/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var adminAuthService *service.AdminAuthService
var adminDonationService *service.DonationService
var adminProofService *service.ProofService
var adminAuditPipeline *service.AuditPipelineService
var adminAuditRepo repository.AuditRepository

func InitAdminHandler(
	authSvc *service.AdminAuthService,
	donations *service.DonationService,
	proofs *service.ProofService,
	pipeline *service.AuditPipelineService,
	audits repository.AuditRepository,
) {
	adminAuthService = authSvc
	adminDonationService = donations
	adminProofService = proofs
	adminAuditPipeline = pipeline
	adminAuditRepo = audits
}

type adminRegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type adminLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminRegister(w http.ResponseWriter, r *http.Request) {
	if adminAuthService == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Admin auth not initialized")
		return
	}
	var body adminRegisterBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Expected JSON body with email and password")
		return
	}
	u, err := adminAuthService.Register(body.Email, body.Password)
	if err != nil {
		status := http.StatusBadRequest
		msg := err.Error()
		if strings.Contains(msg, "E11000") || strings.Contains(strings.ToLower(msg), "duplicate") {
			status = http.StatusConflict
			msg = "An account with this email already exists"
		}
		utils.WriteError(w, status, "REGISTER_FAILED", msg)
		return
	}
	u.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"user": u})
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if adminAuthService == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Admin auth not initialized")
		return
	}
	var body adminLoginBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Expected JSON body with email and password")
		return
	}
	token, u, err := adminAuthService.Login(body.Email, body.Password)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user": map[string]any{
			"id":    u.ID.Hex(),
			"email": u.Email,
		},
	})
}

type adminOverviewProof struct {
	ID             string `json:"id"`
	IPFSHash       string `json:"ipfs_hash"`
	Type           string `json:"type"`
	ExtractedText  string `json:"extracted_text"`
	ExtractionNote string `json:"extraction_note"`
	UploadedAt     int64  `json:"uploaded_at"`
}

type adminOverviewRow struct {
	ID          string               `json:"id"`
	TxHash      string               `json:"tx_hash"`
	DonorWallet string               `json:"donor_wallet"`
	NGOWallet   string               `json:"ngo_wallet"`
	Amount      float64              `json:"amount"`
	Purpose     string               `json:"purpose"`
	Status      string               `json:"status"`
	CreatedAt   int64                `json:"created_at"`
	Proofs      []adminOverviewProof `json:"proofs"`
	LatestAudit *domain.AuditResult  `json:"latest_audit,omitempty"`
}

func AdminOverview(w http.ResponseWriter, r *http.Request) {
	if adminDonationService == nil || adminProofService == nil || adminAuditRepo == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Admin services not initialized")
		return
	}

	donations, err := adminDonationService.ListAllRecent(300)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		return
	}

	out := make([]adminOverviewRow, 0, len(donations))
	for _, d := range donations {
		row := adminOverviewRow{
			ID:          d.ID.Hex(),
			TxHash:      d.TxHash,
			DonorWallet: d.Donor,
			NGOWallet:   d.NGO,
			Amount:      d.Amount,
			Purpose:     d.Purpose,
			Status:      string(d.Status),
			CreatedAt:   d.CreatedAt,
			Proofs:      make([]adminOverviewProof, 0),
		}

		proofs, err := adminProofService.ListByDonationID(d.ID.Hex())
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
			return
		}
		for _, p := range proofs {
			row.Proofs = append(row.Proofs, adminOverviewProof{
				ID:             p.ID.Hex(),
				IPFSHash:       p.IPFSHash,
				Type:           string(p.Type),
				ExtractedText:  p.ExtractedText,
				ExtractionNote: p.ExtractionNote,
				UploadedAt:     p.UploadedAt,
			})
		}

		if ar, err := adminAuditRepo.FindLatestByDonationID(d.ID.Hex()); err == nil && ar != nil {
			row.LatestAudit = ar
		} else if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
			return
		}

		out = append(out, row)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

type adminAuditBody struct {
	DonorProfileLat float64 `json:"donorProfileLat"`
	DonorProfileLng float64 `json:"donorProfileLng"`
	NGOProfileLat   float64 `json:"ngoProfileLat"`
	NGOProfileLng   float64 `json:"ngoProfileLng"`
	RequestLat      float64 `json:"requestLat"`
	RequestLng      float64 `json:"requestLng"`
}

func AdminRunAudit(w http.ResponseWriter, r *http.Request) {
	if adminAuditPipeline == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Audit pipeline not initialized")
		return
	}
	donationID := strings.TrimSpace(chi.URLParam(r, "donationID"))
	if donationID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "donation id required")
		return
	}

	var body adminAuditBody
	_ = json.NewDecoder(r.Body).Decode(&body)

	proofs, err := adminProofService.ListByDonationID(donationID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_DONATION", err.Error())
		return
	}
	if len(proofs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "NO_PROOF", "Upload at least one proof before running an audit")
		return
	}

	var proofLat, proofLng float64
	if len(proofs) > 0 {
		proofLat = proofs[0].Latitude
		proofLng = proofs[0].Longitude
	}

	input := service.AuditInput{
		DonorProfileLat: body.DonorProfileLat,
		DonorProfileLng: body.DonorProfileLng,
		NGOProfileLat:   body.NGOProfileLat,
		NGOProfileLng:   body.NGOProfileLng,
		ProofLat:        proofLat,
		ProofLng:        proofLng,
		RequestLat:      body.RequestLat,
		RequestLng:      body.RequestLng,
	}

	result, err := adminAuditPipeline.RunForDonation(r.Context(), donationID, input)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "AUDIT_FAILED", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"audit": map[string]any{
			"id":         result.ID.Hex(),
			"donation_id": result.DonationID.Hex(),
			"status":     result.Status,
			"score":      result.Score,
			"flags":      result.Flags,
			"created_at": result.CreatedAt,
		},
	})
}

func AdminAuditHistory(w http.ResponseWriter, r *http.Request) {
	if adminAuditRepo == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Audit repository not initialized")
		return
	}

	donationID := strings.TrimSpace(chi.URLParam(r, "donationID"))
	if donationID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "donation id required")
		return
	}

	history, err := adminAuditRepo.FindByDonationID(donationID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "FETCH_FAILED", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(history)
}

func AdminMe(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value(middleware.AdminUserIDKey).(string)
	email, _ := r.Context().Value(middleware.AdminEmailKey).(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id, "email": email})
}
