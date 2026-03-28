package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	coursesAllKey = "courses:all"
	coursesAllTTL = time.Minute * 2
)

type CourseService struct {
	repo       repository.CourseRepo
	lessonRepo repository.LessonRepo
	enrollRepo repository.EnrollmentRepo
	db         *sqlx.DB
	redis      *redis.Client
}

func NewCourseService(
	repo repository.CourseRepo,
	lessonRepo repository.LessonRepo,
	enrollRepo repository.EnrollmentRepo,
	db *sqlx.DB,
	redisClient *redis.Client,
) *CourseService {
	return &CourseService{
		repo:       repo,
		lessonRepo: lessonRepo,
		enrollRepo: enrollRepo,
		db:         db,
		redis:      redisClient,
	}
}

func (cs *CourseService) Create(input models.CreateCourse) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, apperror.BadRequest(err.Error(), err)
	}
	// пока нет у него lessons он не может быть активным
	input.IsActive = false

	id, err := cs.repo.Create(input) // OSYNDA
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

	if err := cs.clearCoursesCache(context.Background()); err != nil {
		slog.Warn("failed to invalidate courses cache after create", "error", err.Error())
	}

	return id, nil
}

func (cs *CourseService) GetAll() ([]models.Course, error) {
	ctx := context.Background()

	if cs.redis != nil {
		cached, err := cs.redis.Get(ctx, coursesAllKey).Result()
		switch {
		case err == nil:
			var cachedCourses []models.Course // '{"id":1,"title":"asdasd"}'
			if uErr := json.Unmarshal([]byte(cached), &cachedCourses); uErr == nil {
				return cachedCourses, nil
			} else {
				slog.Warn("failed to unmarshal courses cache", "error", uErr.Error())
			}
		case errors.Is(err, redis.Nil): // err == в Redis нет такого ключа
			// Значит кеш пустой / тогда возмет данные с БД ниже
		default:
			slog.Warn("failed to get courses from redis", "error", err.Error())
		}
	}

	courses, err := cs.repo.GetAll()
	if err != nil {
		return nil, apperror.Internal("failed to get courses", err)
	}

	if cs.redis != nil {
		jsonCourses, marshalErr := json.Marshal(courses)
		if marshalErr != nil {
			slog.Warn("failed to marshal courses for cache", "error", marshalErr.Error())
		} else if setErr := cs.redis.Set(ctx, coursesAllKey, jsonCourses, coursesAllTTL).Err(); setErr != nil {
			slog.Warn("failed to set courses cache", "error", setErr.Error())
		}
	}

	return courses, nil
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
		return apperror.Internal("failed to start transaction", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err = cs.lessonRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return apperror.Internal("failed to delete related lessons", err)
	}

	if err = cs.enrollRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return apperror.Internal("failed to delete related enrollments", err)
	}

	if err = cs.repo.DeleteByIDTx(ctx, tx, id); err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			return apperror.NotFound("course not found", err)
		}

		return apperror.Internal("failed to delete course", err)
	}

	if err = tx.Commit(); err != nil {
		return apperror.Internal("failed to commit transaction", err)
	}

	if err := cs.clearCoursesCache(context.Background()); err != nil {
		slog.Warn("failed to invalidate courses cache after create", "error", err.Error())
	}

	return nil
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

	if err := cs.clearCoursesCache(ctx); err != nil {
		slog.Warn("failed to invalidate courses cache after update", "error", err.Error())
	}

	return updatedID, nil
}

func (cs *CourseService) clearCoursesCache(ctx context.Context) error {
	if cs.redis == nil {
		return nil
	}

	return cs.redis.Del(ctx, coursesAllKey).Err()
}
