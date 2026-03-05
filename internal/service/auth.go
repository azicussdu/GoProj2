package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azicussdu/GoProj2/internal/auth"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         repository.UserRepo
	tokenManager auth.TokenManager
}

func NewAuthService(userRepo repository.UserRepo, manager auth.TokenManager) *AuthService {
	return &AuthService{
		repo:         userRepo,
		tokenManager: manager,
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

func (s *AuthService) Login(ctx context.Context, input models.LoginUser) (models.AuthTokens, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return models.AuthTokens{}, errors.New("email and password are required")
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.AuthTokens{}, models.ErrInvalidCredentials
		}
		return models.AuthTokens{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return models.AuthTokens{}, models.ErrInvalidCredentials
	}

	accessToken, accessExp, err := s.tokenManager.NewAccessToken(user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	refreshToken, _, err := s.tokenManager.NewRefreshToken(user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	expiresIn := accessExp - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}

	return models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (models.AuthTokens, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return models.AuthTokens{}, errors.New("refresh token is required")
	}

	user, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return models.AuthTokens{}, err
	}

	accessToken, accessExp, err := s.tokenManager.NewAccessToken(*user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	newRefreshToken, _, err := s.tokenManager.NewRefreshToken(*user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	expiresIn := accessExp - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}

	return models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
