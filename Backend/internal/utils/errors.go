package utils

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(APIError{
		Message: message,
		Code:    code,
	})
}
