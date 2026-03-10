package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CourseService struct {
	repo       repository.CourseRepo
	lessonRepo repository.LessonRepo
	enrollRepo repository.EnrollmentRepo
	db         *sqlx.DB
}

func NewCourseService(
	repo repository.CourseRepo,
	lessonRepo repository.LessonRepo,
	enrollRepo repository.EnrollmentRepo,
	db *sqlx.DB,
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
	// пока нет у него lessons он не может быть активным
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
	tx, err := cs.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err = cs.lessonRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return err
	}

	if err = cs.enrollRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return err
	}

	if err = cs.repo.DeleteByIDTx(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (cs *CourseService) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {

	if input.IsActive != nil && *input.IsActive == true {
		lessons, _ := cs.lessonRepo.GetByCourseID(id)
		if len(lessons) == 0 {
			return 0, errors.New("cannot activate course without lessons")
		}
	}

	return cs.repo.Update(ctx, id, input)
}
