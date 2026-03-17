package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type EnrollmentRepo interface {
	Exists(ctx context.Context, userID, courseID int) (bool, error)
	Create(ctx context.Context, input models.CreateEnrollment) (int, error)
	DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error
	DeleteByCourseIDTx(ctx context.Context, tx *gorm.DB, courseID int) error
	GetMyCourses(ctx context.Context, userID int) ([]models.Enrollment, error)
}

type PsgEnrollmentRepo struct {
	db *gorm.DB
}

func NewPsgEnrollmentRepo(db *gorm.DB) *PsgEnrollmentRepo {
	return &PsgEnrollmentRepo{db: db}
}

func (r *PsgEnrollmentRepo) Exists(ctx context.Context, userID, courseID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check enrollment exists: %w", err)
	}
	return count > 0, nil
}

func (r *PsgEnrollmentRepo) Create(ctx context.Context, input models.CreateEnrollment) (int, error) {
	enrollment := models.Enrollment{
		UserID:      input.UserID,
		CourseID:    input.CourseID,
		Progress:    input.Progress,
		IsCompleted: input.IsCompleted,
		CompletedAt: input.CompletedAt,
	}

	tx := r.db.WithContext(ctx).Create(&enrollment)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, models.ErrCourseNotFound
			case "23505":
				return 0, models.ErrEnrollmentAlreadyExists
			}
		}
		return 0, fmt.Errorf("create enrollment: %w", tx.Error)
	}

	return enrollment.ID, nil
}

func (r *PsgEnrollmentRepo) DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error {
	tx := r.db.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).Delete(&models.Enrollment{})
	if tx.Error != nil {
		return fmt.Errorf("delete enrollment: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return models.ErrEnrollmentNotFound
	}
	return nil
}

func (r *PsgEnrollmentRepo) DeleteByCourseIDTx(ctx context.Context, tx *gorm.DB, courseID int) error {
	result := tx.WithContext(ctx).Where("course_id = ?", courseID).Delete(&models.Enrollment{})
	if result.Error != nil {
		return fmt.Errorf("delete enrollments by course: %w", result.Error)
	}
	return nil
}

func (r *PsgEnrollmentRepo) GetMyCourses(ctx context.Context, userID int) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment

	err := r.db.WithContext(ctx).
		Preload("Course").
		Where("user_id = ?", userID).
		Order("enrolled_at DESC").
		Find(&enrollments).Error

	if err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *PsgEnrollmentRepo) GetMyCourses2(ctx context.Context, userID int) ([]models.MyCourse, error) {
	var myCourses []models.MyCourse

	err := r.db.WithContext(ctx).
		Table("enrollments e").
		Select(`
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
		`).
		Joins("JOIN courses c ON c.id = e.course_id").
		Where("e.user_id = ? AND c.deleted_at IS NULL", userID).
		Order("e.enrolled_at DESC").
		Scan(&myCourses).Error
	if err != nil {
		return nil, fmt.Errorf("get my courses: %w", err)
	}

	return myCourses, nil
}
