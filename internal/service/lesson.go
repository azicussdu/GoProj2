package service

import (
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
)

type LessonService struct {
	repo repository.LessonRepo
}

func NewLessonService(repo repository.LessonRepo) *LessonService {
	return &LessonService{repo: repo}
}

func (ls *LessonService) GetAll() ([]models.Lesson, error) {
	return ls.repo.GetAll()
}

func (ls *LessonService) GetByID(id int) (models.Lesson, error) {
	return ls.repo.GetByID(id)
}

func (ls *LessonService) DeleteByID(id int) error {
	return ls.repo.DeleteByID(id)
}

func (ls *LessonService) Create(input models.CreateLesson) (int, error) {
	return ls.repo.Create(input)
}

func (ls *LessonService) Update(id int, input models.UpdateLesson) (int, error) {
	return ls.repo.Update(id, input)
}
