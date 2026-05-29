package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/vishwas-11/trusture-backend/internal/auth"
	"github.com/vishwas-11/trusture-backend/internal/utils"
)

func GetNonce(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	if wallet == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_WALLET", "Wallet is required")
		return
	}

	nonceBytes := make([]byte, 16)
	_, _ = rand.Read(nonceBytes)
	nonce := hex.EncodeToString(nonceBytes)

	auth.SetNonce(wallet, nonce)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"nonce":"` + nonce + `"}`))
}


func VerifyLogin(w http.ResponseWriter, r *http.Request) {
	wallet := r.URL.Query().Get("wallet")
	signature := r.URL.Query().Get("signature")

	if wallet == "" || signature == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Wallet and signature required")
		return
	}

	nonce, ok := auth.GetNonce(wallet)
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_NONCE", "Nonce expired or not found")
		return
	}

	valid, err := auth.VerifySignature(wallet, nonce, signature)
	if err != nil || !valid {
		utils.WriteError(w, http.StatusUnauthorized, "INVALID_SIGNATURE", "Signature verification failed")
		return
	}

	// Important: delete nonce so it cannot be reused
	auth.DeleteNonce(wallet)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Wallet verified"}`))
}
