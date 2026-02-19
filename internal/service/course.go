package service

import (
	"context"
	"errors"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
)

type CourseService struct {
	repo       repository.CourseRepo
	lessonRepo repository.LessonRepo
}

func NewCourseService(repo repository.CourseRepo, lessonRepo repository.LessonRepo) *CourseService {
	return &CourseService{
		repo:       repo,
		lessonRepo: lessonRepo,
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

func (cs *CourseService) DeleteByID(id int) error {
	return cs.repo.DeleteByID(id)
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
