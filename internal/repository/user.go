package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, input models.CreateUser) (int, error)
}

var _ UserRepo = (*PsgUserRepo)(nil)

type PsgUserRepo struct {
	db *sqlx.DB
}

func NewPsgUserRepo(db *sqlx.DB) *PsgUserRepo {
	return &PsgUserRepo{db: db}
}

func (r *PsgUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `
		SELECT
			id,
			full_name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r *PsgUserRepo) Create(ctx context.Context, input models.CreateUser) (int, error) {
	query := `
		INSERT INTO users (
			full_name,
			email,
			password_hash,
			created_at,
			updated_at
		)
		VALUES (
			:full_name,
			:email,
			:password_hash,
			:created_at,
			:updated_at
		)
		RETURNING id
	`

	input.CreatedAt = utils.Now()
	input.UpdatedAt = utils.Now()

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare create user: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.GetContext(ctx, &id, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return 0, models.ErrUserAlreadyExists
			} // email exists
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}
