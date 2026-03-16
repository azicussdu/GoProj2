package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"gorm.io/gorm"
)

type CourseService struct {
	repo       repository.CourseRepo
	lessonRepo repository.LessonRepo
	enrollRepo repository.EnrollmentRepo
	db         *gorm.DB
}

func NewCourseService(
	repo repository.CourseRepo,
	lessonRepo repository.LessonRepo,
	enrollRepo repository.EnrollmentRepo,
	db *gorm.DB,
) *CourseService {
	return &CourseService{
		repo:       repo,
		lessonRepo: lessonRepo,
		enrollRepo: enrollRepo,
		db:         db,
	}
}

func (cs *CourseService) Create(input models.CreateCourse) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}
	input.IsActive = false

	return cs.repo.Create(input)
}

func (cs *CourseService) GetAll() ([]models.Course, error) {
	return cs.repo.GetAll()
}

func (cs *CourseService) GetByID(ctx context.Context, id int) (models.Course, error) {
	return cs.repo.GetByID(ctx, id)
}

func (cs *CourseService) DeleteByID(ctx context.Context, id int) error {
	tx := cs.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := cs.lessonRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := cs.enrollRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := cs.repo.DeleteByIDTx(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (cs *CourseService) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {

	if input.IsActive != nil && *input.IsActive {
		lessons, _ := cs.lessonRepo.GetByCourseID(id)
		if len(lessons) == 0 {
			return 0, errors.New("cannot activate course without lessons")
		}
	}

	return cs.repo.Update(ctx, id, input)
}
