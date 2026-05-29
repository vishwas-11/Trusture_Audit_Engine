package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vishwas-11/trusture-backend/internal/auth"
	"github.com/vishwas-11/trusture-backend/internal/utils"
)

type adminCtxKey string

const (
	AdminUserIDKey adminCtxKey = "adminUserID"
	AdminEmailKey  adminCtxKey = "adminEmail"
)

func AdminBearer(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := strings.TrimSpace(r.Header.Get("Authorization"))
			id, email, err := auth.ParseAdminToken(jwtSecret, raw)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or missing admin token")
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, AdminUserIDKey, id)
			ctx = context.WithValue(ctx, AdminEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
