package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/azicussdu/GoProj2/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

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
		log.Fatalf("connect db error: %v", err)
	}

	password := "admin123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	INSERT INTO users (full_name, email, password_hash, role)
	VALUES ($1, $2, $3, $4)
	`

	_, err = db.ExecContext(
		ctx,
		query,
		"admineke",
		"admin@example.com",
		string(hash),
		"admin",
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin user seeded successfully")
}
