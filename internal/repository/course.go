package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jmoiron/sqlx"
)

type CourseRepo interface {
	GetAll() ([]models.Course, error)
	GetByID(id int) (models.Course, error)
	DeleteByID(id int) error
	Create(input models.CreateCourse) (int, error)
}

type PsgCourseRepo struct {
	db *sqlx.DB
}

func NewPsqCourseRepo(db *sqlx.DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (pcr *PsgCourseRepo) Create(input models.CreateCourse) (int, error) {
	query := `
		INSERT INTO courses (
		    title, description, slug, price, duration, level,
			is_active, instructor_id, created_at, updated_at
		) VALUES (
			:title, :description, :slug, :price, :duration, :level,
			:is_active, :instructor_id, :created_at, :updated_at
		)
		RETURNING id
	`

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	stmt, err := pcr.db.PrepareNamed(query) // Get, Select, Exec
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.Get(&id, input)
	if err != nil {
		return 0, fmt.Errorf("create courses error: %w", err)
	}

	return id, nil
}

func (pcr *PsgCourseRepo) DeleteByID(id int) error {
	query := `
		UPDATE courses
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	result, err := pcr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete course error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (pcr *PsgCourseRepo) GetByID(id int) (models.Course, error) {
	var course models.Course

	query := `
		SELECT id, title, description, slug, price, duration, level,
		level, is_active, instructor_id, created_at, updated_at, deleted_at
		FROM courses
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	err := pcr.db.Get(&course, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Course{}, models.ErrNotFound
		}
		return models.Course{}, fmt.Errorf("get course by id err: %w", err)
	}

	return course, nil
}

func (pcr *PsgCourseRepo) GetAll() ([]models.Course, error) {
	var courses []models.Course

	query := `
		SELECT id, title, description, slug, price, duration, level,
		level, is_active, instructor_id, created_at, updated_at, deleted_at
		FROM courses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := pcr.db.Select(&courses, query)
	if err != nil {
		return nil, err
	}

	return courses, nil
}
