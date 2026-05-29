package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func SignAdminToken(secret string, adminIDHex string, email string, ttl time.Duration) (string, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", errors.New("jwt secret not configured")
	}
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}

	claims := AdminClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   adminIDHex,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseAdminToken(secret string, token string) (adminIDHex string, email string, err error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", "", errors.New("jwt secret not configured")
	}
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	parsed, err := jwt.ParseWithClaims(token, &AdminClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := parsed.Claims.(*AdminClaims)
	if !ok || !parsed.Valid {
		return "", "", errors.New("invalid token")
	}
	return claims.Subject, claims.Email, nil
}
