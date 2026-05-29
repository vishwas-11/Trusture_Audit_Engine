package service

import (
	"errors"
	"strings"
	"time"

	"github.com/vishwas-11/trusture-backend/internal/auth"
	"github.com/vishwas-11/trusture-backend/internal/domain"
	"github.com/vishwas-11/trusture-backend/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AdminAuthService struct {
	admins    repository.AdminRepository
	jwtSecret string
}

func NewAdminAuthService(admins repository.AdminRepository, jwtSecret string) *AdminAuthService {
	return &AdminAuthService{admins: admins, jwtSecret: strings.TrimSpace(jwtSecret)}
}

func (s *AdminAuthService) Register(email string, password string) (*domain.AdminUser, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, errors.New("invalid email")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &domain.AdminUser{
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.admins.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AdminAuthService) Login(email string, password string) (token string, user *domain.AdminUser, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.admins.FindByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token, err = auth.SignAdminToken(s.jwtSecret, u.ID.Hex(), u.Email, 7*24*time.Hour)
	if err != nil {
		return "", nil, err
	}
	u.PasswordHash = ""
	return token, u, nil
}
