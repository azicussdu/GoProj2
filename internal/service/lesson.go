package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
)

type LessonService struct {
	repo       repository.LessonRepo
	courseRepo repository.CourseRepo
}

func NewLessonService(repo repository.LessonRepo, courseRepo repository.CourseRepo) *LessonService {
	return &LessonService{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (ls *LessonService) GetAll() ([]models.Lesson, error) {
	return ls.repo.GetAll()
}

func (ls *LessonService) GetByID(id int) (models.Lesson, error) {
	return ls.repo.GetByID(id)
}

func (ls *LessonService) DeleteByID(ctx context.Context, id int) error {

	lesson, err := ls.repo.GetByID(id)
	if err != nil {
		return err
	}

	course, err := ls.courseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		return err
	}

	if course.IsActive {
		return errors.New("cannot delete lesson inside active course")
	}

	return ls.repo.DeleteByID(ctx, id)
}

func (ls *LessonService) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}

	_, err := ls.courseRepo.GetByID(ctx, input.CourseID)
	if err != nil {
		return 0, err
	}

	return ls.repo.Create(ctx, input)
}

func (ls *LessonService) Update(id int, input models.UpdateLesson) (int, error) {
	return ls.repo.Update(id, input)
}
