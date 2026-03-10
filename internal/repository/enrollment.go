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

type EnrollmentRepo interface {
	Exists(ctx context.Context, userID, courseID int) (bool, error)
	Create(ctx context.Context, input models.CreateEnrollment) (int, error)
	DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error
	DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error
	GetMyCourses(ctx context.Context, userID int) ([]models.MyCourse, error)
}

type PsgEnrollmentRepo struct {
	db *sqlx.DB
}

func NewPsgEnrollmentRepo(db *sqlx.DB) *PsgEnrollmentRepo {
	return &PsgEnrollmentRepo{db: db}
}

func (r *PsgEnrollmentRepo) Exists(ctx context.Context, userID, courseID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM enrollments
			WHERE user_id = $1 AND course_id = $2
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, userID, courseID); err != nil {
		return false, fmt.Errorf("check enrollment exists: %w", err)
	}

	return exists, nil
}

func (r *PsgEnrollmentRepo) Create(ctx context.Context, input models.CreateEnrollment) (int, error) {
	query := `
		INSERT INTO enrollments (
			user_id, course_id, progress, is_completed, enrolled_at, completed_at
		) VALUES (
			:user_id, :course_id, :progress, :is_completed, :enrolled_at, :completed_at
		)
		RETURNING id
	`

	input.EnrolledAt = utils.Now()

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare enrollment insert query: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.GetContext(ctx, &id, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, models.ErrCourseNotFound
			case "23505":
				return 0, models.ErrEnrollmentAlreadyExists
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("create enrollment: %w", err)
	}

	return id, nil
}

func (r *PsgEnrollmentRepo) DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error {
	query := `
		DELETE FROM enrollments
		WHERE user_id = $1 AND course_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, courseID)
	if err != nil {
		return fmt.Errorf("delete enrollment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete enrollment rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrEnrollmentNotFound
	}

	return nil
}

func (r *PsgEnrollmentRepo) DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	query := `
		DELETE FROM enrollments
		WHERE course_id = $1
	`

	if _, err := tx.ExecContext(ctx, query, courseID); err != nil {
		return fmt.Errorf("delete enrollments by course: %w", err)
	}

	return nil
}

func (r *PsgEnrollmentRepo) GetMyCourses(ctx context.Context, userID int) ([]models.MyCourse, error) {
	query := `
		SELECT
			c.id AS course_id,
			c.title,
			c.description,
			c.slug,
			c.price,
			c.duration,
			c.level,
			c.is_active,
			c.teacher_id,
			e.progress,
			e.is_completed,
			e.enrolled_at,
			e.completed_at
		FROM enrollments e
		JOIN courses c ON c.id = e.course_id
		WHERE e.user_id = $1
		  AND c.deleted_at IS NULL
		ORDER BY e.enrolled_at DESC
	`

	var myCourses []models.MyCourse
	if err := r.db.SelectContext(ctx, &myCourses, query, userID); err != nil {
		return nil, fmt.Errorf("get my courses: %w", err)
	}

	return myCourses, nil
}
