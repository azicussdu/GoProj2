package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CourseRepo interface {
	GetAll() ([]models.Course, error)
	GetByID(ctx context.Context, id int) (models.Course, error)
	DeleteByID(id int) error
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error
	Create(input models.CreateCourse) (int, error)
	Update(ctx context.Context, id int, input models.UpdateCourse) (int, error)
}

type PsgCourseRepo struct {
	db *gorm.DB
}

func NewPsqCourseRepo(db *gorm.DB) *PsgCourseRepo {
	return &PsgCourseRepo{db: db}
}

func (pcr *PsgCourseRepo) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {
	updates := map[string]any{}

	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Slug != nil {
		updates["slug"] = *input.Slug
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.Duration != nil {
		updates["duration"] = *input.Duration
	}
	if input.Level != nil {
		updates["level"] = *input.Level
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.TeacherID != nil {
		updates["teacher_id"] = *input.TeacherID
	}

	if len(updates) == 0 {
		return 0, errors.New("no fields to update")
	}

	tx := pcr.db.WithContext(ctx).Model(&models.Course{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, models.ErrTeacherNotFound
			case "23505":
				return 0, models.ErrSlugAlreadyExists
			}
		}
		return 0, fmt.Errorf("update course: %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return 0, models.ErrCourseNotFound
	}

	return id, nil
}

func (pcr *PsgCourseRepo) Create(input models.CreateCourse) (int, error) {
	course := models.Course{
		Title:       input.Title,
		Description: input.Description,
		Slug:        input.Slug,
		Price:       input.Price,
		Duration:    input.Duration,
		Level:       input.Level,
		IsActive:    input.IsActive,
		TeacherID:   input.TeacherID,
	}

	tx := pcr.db.Create(&course)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, models.ErrTeacherNotFound
			case "23505":
				return 0, models.ErrSlugAlreadyExists
			}
		}
		return 0, fmt.Errorf("create course error: %w", tx.Error)
	}

	return course.ID, nil
}

func (pcr *PsgCourseRepo) DeleteByID(id int) error {
	result := pcr.db.Model(&models.Course{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": gorm.Expr("NOW()")})
	if result.Error != nil {
		return fmt.Errorf("delete course error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrCourseNotFound
	}
	return nil
}

func (pcr *PsgCourseRepo) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error {
	result := tx.WithContext(ctx).Model(&models.Course{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": gorm.Expr("NOW()")})
	if result.Error != nil {
		return fmt.Errorf("delete course error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrCourseNotFound
	}
	return nil
}

func (pcr *PsgCourseRepo) GetByID(ctx context.Context, id int) (models.Course, error) {
	var course models.Course
	err := pcr.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Course{}, models.ErrCourseNotFound
		}
		return models.Course{}, fmt.Errorf("get course by id error: %w", err)
	}
	return course, nil
}

func (pcr *PsgCourseRepo) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := pcr.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
