package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type LessonRepo interface {
	GetAll() ([]models.Lesson, error)
	GetByID(id int) (models.Lesson, error)
	GetByCourseID(courseID int) ([]models.Lesson, error)
	DeleteByID(ctx context.Context, id int) error
	Create(ctx context.Context, input models.CreateLesson) (int, error)
	Update(id int, input models.UpdateLesson) (int, error)
}

type PsgLessonRepo struct {
	db *sqlx.DB
}

func NewPsgLessonRepo(db *sqlx.DB) *PsgLessonRepo {
	return &PsgLessonRepo{db: db}
}

func (plr *PsgLessonRepo) GetAll() ([]models.Lesson, error) {
	var lessons []models.Lesson

	query := `
		SELECT id, course_id, title, content, video_url, duration, position,
		is_preview, created_at, updated_at, deleted_at
		FROM lessons
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	if err := plr.db.Select(&lessons, query); err != nil {
		return nil, fmt.Errorf("get all lessons: %w", err)
	}

	return lessons, nil
}

func (plr *PsgLessonRepo) GetByCourseID(courseID int) ([]models.Lesson, error) {
	var lessons []models.Lesson

	query := `
		SELECT id, course_id, title, content, video_url, duration, position,
		is_preview, created_at, updated_at, deleted_at
		FROM lessons
		WHERE course_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	if err := plr.db.Select(&lessons, query); err != nil {
		return nil, fmt.Errorf("get lessons by courseID: %w", err)
	}

	return lessons, nil
}

func (plr *PsgLessonRepo) GetByID(id int) (models.Lesson, error) {
	var lesson models.Lesson

	query := `
		SELECT id, course_id, title, content, video_url, duration, position,
		is_preview, created_at, updated_at, deleted_at
		FROM lessons
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	err := plr.db.Get(&lesson, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Lesson{}, models.ErrLessonNotFound
		}
		return models.Lesson{}, fmt.Errorf("get lesson by id: %w", err)
	}

	return lesson, nil
}

func (plr *PsgLessonRepo) DeleteByID(ctx context.Context, id int) error {
	query := `
		UPDATE lessons
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	result, err := plr.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete lesson: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete lesson rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrLessonNotFound
	}

	return nil
}

func (plr *PsgLessonRepo) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	query := `
		INSERT INTO lessons (
			course_id, title, content, video_url, duration, position,
			is_preview, created_at, updated_at
		) VALUES (
			:course_id, :title, :content, :video_url, :duration, :position,
			:is_preview, :created_at, :updated_at
		)
		RETURNING id
	`

	input.CreatedAt = utils.Now()
	input.UpdatedAt = utils.Now()

	stmt, err := plr.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare lesson insert query: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.GetContext(ctx, &id, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("create lesson: %w", err)
	}

	return id, nil
}

func (plr *PsgLessonRepo) Update(id int, input models.UpdateLesson) (int, error) {
	var setParts []string
	var args []interface{}
	argID := 1

	if input.CourseID != nil {
		setParts = append(setParts, fmt.Sprintf("course_id = $%d", argID))
		args = append(args, *input.CourseID)
		argID++
	}
	if input.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *input.Title)
		argID++
	}
	if input.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argID))
		args = append(args, *input.Content)
		argID++
	}
	if input.VideoURL != nil {
		setParts = append(setParts, fmt.Sprintf("video_url = $%d", argID))
		args = append(args, *input.VideoURL)
		argID++
	}
	if input.Duration != nil {
		setParts = append(setParts, fmt.Sprintf("duration = $%d", argID))
		args = append(args, *input.Duration)
		argID++
	}
	if input.Position != nil {
		setParts = append(setParts, fmt.Sprintf("position = $%d", argID))
		args = append(args, *input.Position)
		argID++
	}
	if input.IsPreview != nil {
		setParts = append(setParts, fmt.Sprintf("is_preview = $%d", argID))
		args = append(args, *input.IsPreview)
		argID++
	}

	if len(setParts) == 0 {
		return 0, errors.New("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, utils.Now())
	argID++

	query := fmt.Sprintf(`
		UPDATE lessons
		SET %s
		WHERE id = $%d
		AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setParts, ", "), argID)
	args = append(args, id)

	var updatedID int
	err := plr.db.Get(&updatedID, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrLessonNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}

		return 0, fmt.Errorf("update lesson: %w", err)
	}

	return updatedID, nil
}
