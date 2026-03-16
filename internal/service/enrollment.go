package service

import (
	"context"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
)

type EnrollmentService struct {
	repo       repository.EnrollmentRepo
	courseRepo repository.CourseRepo
}

func NewEnrollmentService(repo repository.EnrollmentRepo, courseRepo repository.CourseRepo) *EnrollmentService {
	return &EnrollmentService{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (s *EnrollmentService) JoinCourse(ctx context.Context, user models.User, courseID int) (int, error) {
	if user.ID <= 0 {
		return 0, models.ErrUserNotFound
	}

	if user.Role != models.RoleStudent {
		return 0, models.ErrOnlyStudentsCanEnroll
	}

	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		return 0, err
	}

	alreadyEnrolled, err := s.repo.Exists(ctx, user.ID, courseID)
	if err != nil {
		return 0, err
	}
	if alreadyEnrolled {
		return 0, models.ErrEnrollmentAlreadyExists
	}

	input := models.CreateEnrollment{
		UserID:      user.ID,
		CourseID:    courseID,
		Progress:    0,
		IsCompleted: false,
		CompletedAt: nil,
	}

	return s.repo.Create(ctx, input)
}

func (s *EnrollmentService) LeaveCourse(ctx context.Context, user models.User, courseID int) error {
	if user.ID <= 0 {
		return models.ErrUserNotFound
	}

	if user.Role != models.RoleStudent {
		return models.ErrOnlyStudentsCanEnroll
	}

	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		return err
	}

	return s.repo.DeleteByUserAndCourse(ctx, user.ID, courseID)
}

func (s *EnrollmentService) GetMyCourses(ctx context.Context, user models.User) ([]models.MyCourse, error) {
	if user.ID <= 0 {
		return nil, models.ErrUserNotFound
	}

	if user.Role != models.RoleStudent {
		return nil, models.ErrOnlyStudentsCanEnroll
	}

	return s.repo.GetMyCourses(ctx, user.ID)
}
