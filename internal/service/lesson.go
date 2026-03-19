package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/jmoiron/sqlx"
)

type LessonService struct {
	repo       repository.LessonRepo
	courseRepo repository.CourseRepo
	db         *sqlx.DB
}

func NewLessonService(repo repository.LessonRepo, courseRepo repository.CourseRepo, db *sqlx.DB) *LessonService {
	return &LessonService{
		repo:       repo,
		courseRepo: courseRepo,
		db:         db,
	}
}

func (ls *LessonService) GetAll() ([]models.Lesson, error) {
	lessons, err := ls.repo.GetAll()
	if err != nil {
		return nil, apperror.Internal("failed to get lessons", err)
	}

	return lessons, nil
}

func (ls *LessonService) GetByID(id int) (models.Lesson, error) {
	lesson, err := ls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			return models.Lesson{}, apperror.NotFound("lesson not found", err)
		}

		return models.Lesson{}, apperror.Internal("failed to get lesson", err)
	}

	return lesson, nil
}

func (ls *LessonService) DeleteByID(ctx context.Context, id int) error {

	lesson, err := ls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			return apperror.NotFound("lesson not found", err)
		}

		return apperror.Internal("failed to get lesson", err)
	}

	course, err := ls.courseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return apperror.NotFound("course not found", err)
		}

		return apperror.Internal("failed to get course", err)
	}

	if course.IsActive {
		return apperror.Conflict("cannot delete lesson inside active course", errors.New("cannot delete lesson inside active course"))
	}

	err = ls.repo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			return apperror.NotFound("lesson not found", err)
		}

		return apperror.Internal("failed to delete lesson", err)
	}

	return nil
}

func (ls *LessonService) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, apperror.BadRequest(err.Error(), err)
	}

	_, err := ls.courseRepo.GetByID(ctx, input.CourseID)
	if err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return 0, apperror.NotFound("course not found", err)
		}

		return 0, apperror.Internal("failed to get course", err)
	}

	id, err := ls.repo.Create(ctx, input)
	if err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return 0, apperror.NotFound("course not found", err)
		}

		return 0, apperror.Internal("failed to create lesson", err)
	}

	return id, nil
}

func (ls *LessonService) Update(id int, input models.UpdateLesson) (int, error) {
	updatedID, err := ls.repo.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrLessonNotFound):
			return 0, apperror.NotFound("lesson not found", err)
		case errors.Is(err, models.ErrCourseNotFound):
			return 0, apperror.NotFound("course not found", err)
		default:
			return 0, apperror.Internal("failed to update lesson", err)
		}
	}

	return updatedID, nil
}
