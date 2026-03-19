package service

import (
	"context"
	"errors"
	"strings"

	"github.com/azicussdu/GoProj2/internal/apperror"
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
		return 0, apperror.Unauthorized("user is not authenticated", models.ErrUserNotFound)
	}

	role := strings.TrimSpace(strings.ToLower(user.Role))
	if role != models.RoleStudent {
		return 0, apperror.Forbidden(models.ErrOnlyStudentsCanEnroll.Error(), models.ErrOnlyStudentsCanEnroll)
	}

	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return 0, apperror.NotFound("course not found", err)
		}

		return 0, apperror.Internal("failed to get course", err)
	}

	alreadyEnrolled, err := s.repo.Exists(ctx, user.ID, courseID)
	if err != nil {
		return 0, apperror.Internal("failed to check enrollment", err)
	}
	if alreadyEnrolled {
		return 0, apperror.Conflict(models.ErrEnrollmentAlreadyExists.Error(), models.ErrEnrollmentAlreadyExists)
	}

	input := models.CreateEnrollment{
		UserID:      user.ID,
		CourseID:    courseID,
		Progress:    0,
		IsCompleted: false,
		CompletedAt: nil,
	}

	id, err := s.repo.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCourseNotFound):
			return 0, apperror.NotFound("course not found", err)
		case errors.Is(err, models.ErrEnrollmentAlreadyExists):
			return 0, apperror.Conflict(models.ErrEnrollmentAlreadyExists.Error(), err)
		default:
			return 0, apperror.Internal("failed to enroll in course", err)
		}
	}

	return id, nil
}

func (s *EnrollmentService) LeaveCourse(ctx context.Context, user models.User, courseID int) error {
	if user.ID <= 0 {
		return apperror.Unauthorized("user is not authenticated", models.ErrUserNotFound)
	}

	role := strings.TrimSpace(strings.ToLower(user.Role))
	if role != models.RoleStudent {
		return apperror.Forbidden(models.ErrOnlyStudentsCanEnroll.Error(), models.ErrOnlyStudentsCanEnroll)
	}

	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return apperror.NotFound("course not found", err)
		}

		return apperror.Internal("failed to get course", err)
	}

	if err := s.repo.DeleteByUserAndCourse(ctx, user.ID, courseID); err != nil {
		if errors.Is(err, models.ErrEnrollmentNotFound) {
			return apperror.NotFound(models.ErrEnrollmentNotFound.Error(), err)
		}

		return apperror.Internal("failed to leave course", err)
	}

	return nil
}

func (s *EnrollmentService) GetMyCourses(ctx context.Context, user models.User) ([]models.MyCourse, error) {
	if user.ID <= 0 {
		return nil, apperror.Unauthorized("user is not authenticated", models.ErrUserNotFound)
	}

	role := strings.TrimSpace(strings.ToLower(user.Role))
	if role != models.RoleStudent {
		return nil, apperror.Forbidden(models.ErrOnlyStudentsCanEnroll.Error(), models.ErrOnlyStudentsCanEnroll)
	}

	courses, err := s.repo.GetMyCourses(ctx, user.ID)
	if err != nil {
		return nil, apperror.Internal("failed to get my courses", err)
	}

	return courses, nil
}
