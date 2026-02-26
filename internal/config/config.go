package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	LogLevel string

	Database *DBConfig
	JWT      *JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	accessTTL, err := parseDurationEnv("JWT_ACCESS_TTL", "15m")
	if err != nil {
		return nil, err
	}

	refreshTTL, err := parseDurationEnv("JWT_REFRESH_TTL", "720h")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "INFO"),

		Database: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("SSL_MODE", "disable"),
		},
		JWT: &JWTConfig{
			Secret:     getEnv("JWT_SECRET", "secret-key"),
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
			Issuer:     getEnv("JWT_ISSUER", "company-name"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func parseDurationEnv(key, defaultValue string) (time.Duration, error) {
	value := getEnv(key, defaultValue) // value = "30m"
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return duration, nil
}
