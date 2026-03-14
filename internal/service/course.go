package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/apperror"
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

	id, err := cs.repo.Create(input)
	if err != nil {

		switch {
		case errors.Is(err, models.ErrTeacherNotFound):
			return 0, apperror.NotFound("teacher not found", err)

		case errors.Is(err, models.ErrSlugAlreadyExists):
			return 0, apperror.Conflict("slug already exists", err)

		default:
			return 0, apperror.Internal("failed to create course", err)
		}
	}

	return id, nil
}

func (cs *CourseService) GetAll() ([]models.Course, error) {
	return cs.repo.GetAll()
}

func (cs *CourseService) GetByID(ctx context.Context, id int) (models.Course, error) {
	course, err := cs.repo.GetByID(ctx, id)
	if err != nil {

		if errors.Is(err, models.ErrCourseNotFound) {
			return models.Course{}, apperror.NotFound("course not found", err)
		}

		return models.Course{}, apperror.Internal("failed to get course", err)
	}

	return course, nil
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
		if errors.Is(err, models.ErrCourseNotFound) {
			return apperror.NotFound("course not found", err)
		}

		return apperror.Internal("failed to delete course", err)
	}

	return tx.Commit()
}

func (cs *CourseService) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {

	if input.IsActive != nil && *input.IsActive {

		lessons, err := cs.lessonRepo.GetByCourseID(id)
		if err != nil {
			return 0, apperror.Internal("failed to check lessons", err)
		}

		if len(lessons) == 0 {
			return 0, apperror.BadRequest(models.ErrCourseCannotBeActivated.Error(), models.ErrCourseCannotBeActivated)
		}
	}

	updatedID, err := cs.repo.Update(ctx, id, input)
	if err != nil {

		if errors.Is(err, models.ErrCourseNotFound) {
			return 0, apperror.NotFound("course not found", err)
		}

		if errors.Is(err, models.ErrSlugAlreadyExists) {
			return 0, apperror.Conflict("slug already exists", err)
		}

		return 0, apperror.Internal("failed to update course", err)
	}

	return updatedID, nil
}
