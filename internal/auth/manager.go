package auth

import "github.com/azicussdu/GoProj2/internal/models"

type TokenManager interface {
	NewAccessToken(user models.User) (string, int64, error)
}
