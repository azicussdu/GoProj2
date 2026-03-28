package service

import (
	"context"

	"github.com/azicussdu/GoProj2/internal/models"
)

type CourseServiceI interface {
	Create(input models.CreateCourse) (int, error)
	GetAll() ([]models.Course, error)
	GetByID(ctx context.Context, id int) (models.Course, error)
	DeleteByID(ctx context.Context, id int) error
	Update(ctx context.Context, id int, input models.UpdateCourse) (int, error)
}

type Services struct {
	Course     CourseServiceI
	Lesson     *LessonService
	Enrollment *EnrollmentService
	Auth       *AuthService
}
