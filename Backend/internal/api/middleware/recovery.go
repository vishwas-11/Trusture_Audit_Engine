package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/vishwas-11/trusture-backend/internal/utils"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				utils.WriteError(
					w,
					http.StatusInternalServerError,
					"PANIC",
					"Internal server error",
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
