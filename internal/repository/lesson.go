package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type LessonRepo interface {
	GetAll() ([]models.Lesson, error)
	GetByID(id int) (models.Lesson, error)
	GetByCourseID(courseID int) ([]models.Lesson, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteByCourseIDTx(ctx context.Context, tx *gorm.DB, courseID int) error
	Create(ctx context.Context, input models.CreateLesson) (int, error)
	Update(id int, input models.UpdateLesson) (int, error)
}

type PsgLessonRepo struct {
	db *gorm.DB
}

func NewPsgLessonRepo(db *gorm.DB) *PsgLessonRepo {
	return &PsgLessonRepo{db: db}
}

func (plr *PsgLessonRepo) GetAll() ([]models.Lesson, error) {
	var lessons []models.Lesson
	if err := plr.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("get all lessons: %w", err)
	}
	return lessons, nil
}

func (plr *PsgLessonRepo) GetByCourseID(courseID int) ([]models.Lesson, error) {
	var lessons []models.Lesson
	if err := plr.db.Where("course_id = ? AND deleted_at IS NULL", courseID).Order("created_at ASC").Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("get lessons by courseID: %w", err)
	}
	return lessons, nil
}

func (plr *PsgLessonRepo) GetByID(id int) (models.Lesson, error) {
	var lesson models.Lesson
	err := plr.db.Where("id = ? AND deleted_at IS NULL", id).First(&lesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Lesson{}, models.ErrLessonNotFound
		}
		return models.Lesson{}, fmt.Errorf("get lesson by id: %w", err)
	}
	return lesson, nil
}

func (plr *PsgLessonRepo) DeleteByID(ctx context.Context, id int) error {
	result := plr.db.WithContext(ctx).Model(&models.Lesson{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": gorm.Expr("NOW()")})
	if result.Error != nil {
		return fmt.Errorf("delete lesson: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrLessonNotFound
	}
	return nil
}

func (plr *PsgLessonRepo) DeleteByCourseIDTx(ctx context.Context, tx *gorm.DB, courseID int) error {
	result := tx.WithContext(ctx).Model(&models.Lesson{}).
		Where("course_id = ? AND deleted_at IS NULL", courseID).
		Updates(map[string]any{"deleted_at": gorm.Expr("NOW()")})
	if result.Error != nil {
		return fmt.Errorf("delete lessons by course: %w", result.Error)
	}
	return nil
}

func (plr *PsgLessonRepo) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	lesson := models.Lesson{
		CourseID:  input.CourseID,
		Title:     input.Title,
		Content:   input.Content,
		VideoURL:  input.VideoURL,
		Duration:  input.Duration,
		Position:  input.Position,
		IsPreview: input.IsPreview,
	}

	tx := plr.db.WithContext(ctx).Create(&lesson)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("create lesson: %w", tx.Error)
	}

	return lesson.ID, nil
}

func (plr *PsgLessonRepo) Update(id int, input models.UpdateLesson) (int, error) {
	updates := map[string]any{}

	if input.CourseID != nil {
		updates["course_id"] = *input.CourseID
	}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.VideoURL != nil {
		updates["video_url"] = *input.VideoURL
	}
	if input.Duration != nil {
		updates["duration"] = *input.Duration
	}
	if input.Position != nil {
		updates["position"] = *input.Position
	}
	if input.IsPreview != nil {
		updates["is_preview"] = *input.IsPreview
	}

	if len(updates) == 0 {
		return 0, errors.New("no fields to update")
	}

	tx := plr.db.Model(&models.Lesson{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("update lesson: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return 0, models.ErrLessonNotFound
	}

	return id, nil
}
