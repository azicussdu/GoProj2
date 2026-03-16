package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, input models.CreateUser) (int, error)
	GetByID(id int) (models.User, error)
	UpdateRole(id int, role string) (int, error)
}

var _ UserRepo = (*PsgUserRepo)(nil)

type PsgUserRepo struct {
	db *gorm.DB
}

func NewPsgUserRepo(db *gorm.DB) *PsgUserRepo {
	return &PsgUserRepo{db: db}
}

func (r *PsgUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *PsgUserRepo) Create(ctx context.Context, input models.CreateUser) (int, error) {
	user := models.User{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	}

	tx := r.db.WithContext(ctx).Create(&user)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) && pgErr.Code == "23505" {
			return 0, models.ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("create user: %w", tx.Error)
	}

	return user.ID, nil
}

func (r *PsgUserRepo) GetByID(id int) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *PsgUserRepo) UpdateRole(userID int, role string) (int, error) {
	tx := r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role)
	if tx.Error != nil {
		return 0, fmt.Errorf("update user role: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return 0, models.ErrRoleChangeOnlyFromStudent
	}
	return userID, nil
}
