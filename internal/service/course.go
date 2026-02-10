package service

import (
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
)

type CourseService struct {
	repo repository.CourseRepo // interface
}

func NewCourseService(repo repository.CourseRepo) *CourseService {
	return &CourseService{repo: repo}
}

func (cs *CourseService) Create(input models.CreateCourse) (int, error) {
	// will add some logic
	return cs.repo.Create(input)
}

func (cs *CourseService) GetAll() ([]models.Course, error) {
	return cs.repo.GetAll()
}

func (cs *CourseService) GetByID(id int) (models.Course, error) {
	return cs.repo.GetByID(id)
}

func (cs *CourseService) DeleteByID(id int) error {
	return cs.repo.DeleteByID(id)
}
