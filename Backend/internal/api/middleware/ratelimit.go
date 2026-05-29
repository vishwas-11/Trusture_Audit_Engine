package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vishwas-11/trusture-backend/internal/utils"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

const (
	maxRequests = 100
	window      = time.Minute
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > window {
			v.count = 1
			v.lastSeen = time.Now()
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if v.count >= maxRequests {
			mu.Unlock()
			utils.WriteError(
				w,
				http.StatusTooManyRequests,
				"RATE_LIMITED",
				"Too many requests, slow down",
			)
			return
		}

		v.count++
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
