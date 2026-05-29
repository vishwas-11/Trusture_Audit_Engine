package handlers

import (
	"net/http"
)

func DonorDashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Donor dashboard data"))
}
