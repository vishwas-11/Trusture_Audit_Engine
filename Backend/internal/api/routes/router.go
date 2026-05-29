package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vishwas-11/trusture-backend/internal/api/handlers"
	"github.com/vishwas-11/trusture-backend/internal/api/middleware"
	"github.com/vishwas-11/trusture-backend/internal/config"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit)

	// --- Public routes ---
	r.Get("/auth/nonce", handlers.GetNonce)
	r.Get("/auth/verify", handlers.VerifyLogin)

	// --- Health check ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// --- Donor routes ---
	r.Route("/donor", func(r chi.Router) {
		r.Use(middleware.RoleRequired("DONOR")) // TEMP (replaced in Task 3)
		r.Get("/dashboard", handlers.DonorDashboard)
	})

	// --- NGO routes ---
	r.Route("/ngo", func(r chi.Router) {
		r.Use(middleware.RoleRequired("NGO")) // TEMP (replaced in Task 3)
		r.Get("/dashboard", handlers.NGODashboard)
		r.Post("/proof", handlers.UploadProof)
	})

	// --- Admin routes (email/password + JWT) ---
	r.Route("/admin", func(r chi.Router) {
		r.Post("/auth/register", handlers.AdminRegister)
		r.Post("/auth/login", handlers.AdminLogin)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminBearer(cfg.JWTSecret))
			r.Get("/me", handlers.AdminMe)
			r.Get("/overview", handlers.AdminOverview)
			r.Post("/donations/{donationID}/audit", handlers.AdminRunAudit)
			r.Get("/donations/{donationID}/audits", handlers.AdminAuditHistory)
		})
	})

	return r
}
