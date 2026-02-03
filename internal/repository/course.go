package repository

import "github.com/azicussdu/GoProj2/internal/models"

type CourseRepo interface {
	GetAll() ([]models.Course, error)
	// TODO реализуй остальные методы
}

type PsgCourseRepo struct {
	db *DB
}

func NewPsqCourseRepo(db *DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (pcr *PsgCourseRepo) GetAll() ([]models.Course, error) {
	// TODO db koldanu kerek (chtovy posgtgrestan dannyi tartu ushin)
	return []models.Course{
		{ID: 1, Name: "Go Basics"},
		{ID: 2, Name: "Nodejs Basics"},
		{ID: 3, Name: "Java Basics"},
	}, nil
}
