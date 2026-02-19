package service

import (
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

func (ls *LessonService) DeleteByID(id int) error {

	lesson, err := ls.repo.GetByID(id)
	if err != nil {
		return err
	}

	course, err := ls.courseRepo.GetByID(lesson.CourseID)
	if err != nil {
		return err
	}

	if course.IsActive {
		return errors.New("cannot delete lesson without active lessons")
	}

	return ls.repo.DeleteByID(id)
}

func (ls *LessonService) Create(input models.CreateLesson) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}

	_, err := ls.courseRepo.GetByID(input.CourseID)
	if err != nil {
		return 0, err
	}

	return ls.repo.Create(input)
}

func (ls *LessonService) Update(id int, input models.UpdateLesson) (int, error) {
	return ls.repo.Update(id, input)
}
