package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/service"
	"github.com/vishwas-11/trusture-backend/internal/utils"
)

var proofService *service.ProofService
var documentExtractService *service.DocumentExtractService

func InitProofHandler(
	proofs *service.ProofService,
	extract *service.DocumentExtractService,
) {
	proofService = proofs
	documentExtractService = extract
}

func UploadProof(w http.ResponseWriter, r *http.Request) {
	if proofService == nil || documentExtractService == nil {
		utils.WriteError(w, http.StatusInternalServerError, "SERVICE_NOT_READY", "Proof services are not initialized")
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

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_MULTIPART",
			"Invalid or too large multipart form",
		)
		return
	}

	donationID := strings.TrimSpace(r.FormValue("donationId"))
	if donationID == "" {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_DONATION",
			"Donation ID is required",
		)
		return
	}

	file, header, err := r.FormFile("proof")
	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"MISSING_FILE",
			"Proof file is required",
		)
		return
	}
	defer file.Close()

	origName := header.Filename
	if strings.TrimSpace(origName) == "" {
		origName = "proof"
	}

	tmpFile, err := os.CreateTemp("", "proof-*")
	if err != nil {
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"FILE_ERROR",
			"Failed to create temporary file",
		)
		return
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.ReadFrom(file); err != nil {
		_ = tmpFile.Close()
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"FILE_WRITE_FAILED",
			"Failed to write uploaded file",
		)
		return
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		utils.WriteError(w, http.StatusInternalServerError, "FILE_SYNC_FAILED", err.Error())
		return
	}
	if err := tmpFile.Close(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "FILE_CLOSE_FAILED", err.Error())
		return
	}

	ipfs := service.NewPinataService()
	hash, err := ipfs.UploadFile(tmpPath)
	if err != nil {
		log.Printf("❌ Pinata upload failed for donation=%s wallet=%s: %v", donationID, wallet, err)
		utils.WriteError(
			w,
			http.StatusInternalServerError,
			"IPFS_UPLOAD_FAILED",
			"Failed to upload proof to IPFS: "+err.Error(),
		)
		return
	}

	extractedText, extractionNote := documentExtractService.ExtractFromFile(tmpPath, origName)

	proofType := domain.ProofType(strings.ToUpper(strings.TrimSpace(r.FormValue("proofType"))))
	if proofType == "" {
		proofType = domain.ProofMedia
	}
	geoTag := strings.TrimSpace(r.FormValue("geoTag"))
	latitude := parseFloat(r.FormValue("proofLat"))
	longitude := parseFloat(r.FormValue("proofLng"))
	gstin := strings.ToUpper(strings.TrimSpace(r.FormValue("gstin")))
	vendorName := strings.TrimSpace(r.FormValue("vendorName"))

	proof, err := proofService.CreateProof(
		donationID,
		hash,
		proofType,
		geoTag,
		latitude,
		longitude,
		gstin,
		vendorName,
		extractedText,
		extractionNote,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "PROOF_SAVE_FAILED", "Failed to save proof metadata")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ipfsHash":        hash,
		"proofId":         proof.ID.Hex(),
		"extractedText":   extractedText,
		"extractionNote":  extractionNote,
		"auditPending":    true,
		"message":         "Proof uploaded. An administrator must run the audit from the admin dashboard when ready.",
	})
}

func parseFloat(raw string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0
	}
	return value
}
