package service

import (
	"context"
	"errors"
	"strings"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepo
	//tokenManager TokenManager
}

func NewAuthService(userRepo repository.UserRepo) *AuthService {
	return &AuthService{
		repo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, input models.RegisterUser) (int, error) {
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.FullName == "" || input.Email == "" || input.Password == "" {
		return 0, errors.New("full_name, email, and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := models.CreateUser{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
