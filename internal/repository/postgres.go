package repository

import "github.com/azicussdu/GoProj2/internal/config"

// fake DB (later sqlx.DB)
type DB struct {
	connectionPath string
}

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	db := &DB{}
	return db, nil
}
