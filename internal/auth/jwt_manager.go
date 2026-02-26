package auth

import (
	"strconv"
	"time"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int    `json:"uid"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
}

func (jm *JWTManager) NewAccessToken(user models.User) (string, int64, error) {
	expiresAt := time.Now().Add(jm.accessTTL)

	claims := Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jm.issuer,
			Subject:   strconv.Itoa(user.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // number in seconds
			ExpiresAt: jwt.NewNumericDate(expiresAt),  // number in seconds
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(jm.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, expiresAt.Unix(), nil

}
