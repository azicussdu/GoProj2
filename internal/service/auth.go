package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azicussdu/GoProj2/internal/apperror"
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
		return 0, apperror.BadRequest("full_name, email, and password are required", errors.New("full_name, email, and password are required"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, apperror.Internal("failed to hash password", err)
	}

	user := models.CreateUser{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return 0, apperror.Conflict("user already exists", err)
		}

		return 0, apperror.Internal("failed to register user", err)
	}

	return id, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginUser) (models.AuthTokens, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return models.AuthTokens{}, apperror.BadRequest("email and password are required", errors.New("email and password are required"))
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.AuthTokens{}, apperror.Unauthorized("invalid email or password", models.ErrInvalidCredentials)
		}
		return models.AuthTokens{}, apperror.Internal("failed to get user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return models.AuthTokens{}, apperror.Unauthorized("invalid email or password", models.ErrInvalidCredentials)
	}

	accessToken, accessExp, err := s.tokenManager.NewAccessToken(user)
	if err != nil {
		return models.AuthTokens{}, apperror.Internal("failed to create access token", err)
	}

	refreshToken, _, err := s.tokenManager.NewRefreshToken(user)
	if err != nil {
		return models.AuthTokens{}, apperror.Internal("failed to create refresh token", err)
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
		return models.AuthTokens{}, apperror.BadRequest("refresh token is required", errors.New("refresh token is required"))
	}

	user, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return models.AuthTokens{}, apperror.Unauthorized("invalid refresh token", err)
	}

	accessToken, accessExp, err := s.tokenManager.NewAccessToken(*user)
	if err != nil {
		return models.AuthTokens{}, apperror.Internal("failed to create access token", err)
	}

	newRefreshToken, _, err := s.tokenManager.NewRefreshToken(*user)
	if err != nil {
		return models.AuthTokens{}, apperror.Internal("failed to create refresh token", err)
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

func (s *AuthService) ChangeUserRole(userID int, newRole string) (int, error) {
	role := strings.TrimSpace(strings.ToLower(newRole))

	if !models.IsValidRole(role) {
		return 0, apperror.BadRequest("role must be teacher or admin", models.ErrInvalidRole)
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return 0, apperror.NotFound("user not found", err)
		}

		return 0, apperror.Internal("failed to get user", err)
	}

	if strings.ToLower(strings.TrimSpace(user.Role)) != models.RoleStudent {
		return 0, apperror.Conflict("only student role can be changed", models.ErrRoleChangeOnlyFromStudent)
	}

	updatedID, err := s.repo.UpdateRole(userID, role)
	if err != nil {
		return 0, apperror.Internal("failed to update user role", err)
	}

	return updatedID, nil
}
