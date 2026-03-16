package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int             `json:"uid"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	TokenType string          `json:"type"`
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

func (jm *JWTManager) ParseAccessToken(tokenStr string) (*models.User, error) {
	claims, err := jm.parseToken(tokenStr, "access")
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

func (jm *JWTManager) NewRefreshToken(user models.User) (string, int64, error) {
	expiresAt := time.Now().Add(jm.refreshTTL)

	claims := Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jm.issuer,
			Subject:   strconv.Itoa(user.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(jm.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, expiresAt.Unix(), nil
}

func (m *JWTManager) ParseRefreshToken(tokenStr string) (*models.User, error) {
	claims, err := m.parseToken(tokenStr, "refresh")
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

func (jm *JWTManager) parseToken(tokenStr, expectedType string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is required")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	parsedToken, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jm.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("unexpected token type: %s", claims.TokenType)
	}

	if jm.issuer != "" && claims.Issuer != jm.issuer {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}
