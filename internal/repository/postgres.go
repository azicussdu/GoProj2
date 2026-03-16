package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/azicussdu/GoProj2/internal/config"
	"github.com/azicussdu/GoProj2/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect db error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if cfg.Database.AutoMigrate {
		err = db.AutoMigrate(
			&models.User{},
			&models.Course{},
			&models.Lesson{},
			&models.Enrollment{},
			&models.LessonCompletion{},
		)
		if err != nil {
			return nil, fmt.Errorf("auto-migrate db error: %w", err)
		}
		slog.Info("GORM auto-migration completed")
	}

	slog.Info("PostgreSQL connected successfully")

	return db, nil
}
