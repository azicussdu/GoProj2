package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/azicussdu/GoProj2/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	slog.Info("PostgreSQL connected successfully")

	return db, nil
}
