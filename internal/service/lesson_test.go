package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLessonService_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(repo *MockLessonRepo)
		expected   []models.Lesson
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "success",
			setup: func(repo *MockLessonRepo) {
				repo.On("GetAll").Return([]models.Lesson{
					{ID: 1, CourseID: 10, Title: "Intro"},
					{ID: 2, CourseID: 10, Title: "Advanced"},
				}, nil).Once()
			},
			expected: []models.Lesson{
				{ID: 1, CourseID: 10, Title: "Intro"},
				{ID: 2, CourseID: 10, Title: "Advanced"},
			},
		},
		{
			name: "internal error",
			setup: func(repo *MockLessonRepo) {
				repo.On("GetAll").Return(nil, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get lessons",
			errWrapped: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := new(MockLessonRepo)
			tt.setup(lessonRepo)
			svc := NewLessonService(lessonRepo, new(MockCourseRepo), nil)

			lessons, err := svc.GetAll()

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, lessons)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			lessonRepo.AssertExpectations(t)
		})
	}
}

func TestLessonService_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		setup      func(repo *MockLessonRepo)
		expected   models.Lesson
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "success",
			id:   1,
			setup: func(repo *MockLessonRepo) {
				repo.On("GetByID", 1).Return(models.Lesson{ID: 1, CourseID: 2, Title: "Lesson"}, nil).Once()
			},
			expected: models.Lesson{ID: 1, CourseID: 2, Title: "Lesson"},
		},
		{
			name: "not found",
			id:   2,
			setup: func(repo *MockLessonRepo) {
				repo.On("GetByID", 2).Return(models.Lesson{}, models.ErrLessonNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "lesson not found",
			errWrapped: models.ErrLessonNotFound,
		},
		{
			name: "internal error",
			id:   3,
			setup: func(repo *MockLessonRepo) {
				repo.On("GetByID", 3).Return(models.Lesson{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get lesson",
			errWrapped: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := new(MockLessonRepo)
			tt.setup(lessonRepo)
			svc := NewLessonService(lessonRepo, new(MockCourseRepo), nil)

			lesson, err := svc.GetByID(tt.id)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, lesson)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			lessonRepo.AssertExpectations(t)
		})
	}
}

func TestLessonService_Create(t *testing.T) {
	ctx := context.Background()
	content := "body"

	tests := []struct {
		name       string
		input      models.CreateLesson
		setup      func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo)
		expectedID int
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "bad request validation",
			input: models.CreateLesson{
				CourseID: 0,
				Title:    "",
			},
			errCode: http.StatusBadRequest,
			errMsg:  "invalid course id",
		},
		{
			name: "course not found in course repo",
			input: models.CreateLesson{
				CourseID:  7,
				Title:     "Lesson",
				Content:   &content,
				Duration:  10,
				Position:  1,
				IsPreview: true,
			},
			setup: func(_ *MockLessonRepo, courseRepo *MockCourseRepo) {
				courseRepo.On("GetByID", ctx, 7).Return(models.Course{}, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal error in course repo",
			input: models.CreateLesson{
				CourseID:  8,
				Title:     "Lesson",
				Content:   &content,
				Duration:  10,
				Position:  1,
				IsPreview: true,
			},
			setup: func(_ *MockLessonRepo, courseRepo *MockCourseRepo) {
				courseRepo.On("GetByID", ctx, 8).Return(models.Course{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get course",
			errWrapped: errDB,
		},
		{
			name: "course not found in lesson repo create",
			input: models.CreateLesson{
				CourseID:  9,
				Title:     "Lesson",
				Content:   &content,
				Duration:  10,
				Position:  1,
				IsPreview: true,
			},
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				courseRepo.On("GetByID", ctx, 9).Return(models.Course{ID: 9}, nil).Once()
				lessonRepo.On("Create", ctx, models.CreateLesson{
					CourseID:  9,
					Title:     "Lesson",
					Content:   &content,
					Duration:  10,
					Position:  1,
					IsPreview: true,
				}).Return(0, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal error in lesson repo create",
			input: models.CreateLesson{
				CourseID:  10,
				Title:     "Lesson",
				Content:   &content,
				Duration:  10,
				Position:  1,
				IsPreview: true,
			},
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				courseRepo.On("GetByID", ctx, 10).Return(models.Course{ID: 10}, nil).Once()
				lessonRepo.On("Create", ctx, models.CreateLesson{
					CourseID:  10,
					Title:     "Lesson",
					Content:   &content,
					Duration:  10,
					Position:  1,
					IsPreview: true,
				}).Return(0, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to create lesson",
			errWrapped: errDB,
		},
		{
			name: "success",
			input: models.CreateLesson{
				CourseID:  11,
				Title:     "Lesson",
				Content:   &content,
				Duration:  10,
				Position:  1,
				IsPreview: false,
			},
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				courseRepo.On("GetByID", ctx, 11).Return(models.Course{ID: 11}, nil).Once()
				lessonRepo.On("Create", ctx, models.CreateLesson{
					CourseID:  11,
					Title:     "Lesson",
					Content:   &content,
					Duration:  10,
					Position:  1,
					IsPreview: false,
				}).Return(99, nil).Once()
			},
			expectedID: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := new(MockLessonRepo)
			courseRepo := new(MockCourseRepo)
			if tt.setup != nil {
				tt.setup(lessonRepo, courseRepo)
			}

			svc := NewLessonService(lessonRepo, courseRepo, nil)
			id, err := svc.Create(ctx, tt.input)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			lessonRepo.AssertExpectations(t)
			courseRepo.AssertExpectations(t)
		})
	}
}

func TestLessonService_Update(t *testing.T) {
	title := "Updated lesson"
	input := models.UpdateLesson{Title: &title}

	tests := []struct {
		name       string
		setup      func(repo *MockLessonRepo)
		expectedID int
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "lesson not found",
			setup: func(repo *MockLessonRepo) {
				repo.On("Update", 1, input).Return(0, models.ErrLessonNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "lesson not found",
			errWrapped: models.ErrLessonNotFound,
		},
		{
			name: "course not found",
			setup: func(repo *MockLessonRepo) {
				repo.On("Update", 1, input).Return(0, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal error",
			setup: func(repo *MockLessonRepo) {
				repo.On("Update", 1, input).Return(0, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to update lesson",
			errWrapped: errDB,
		},
		{
			name: "success",
			setup: func(repo *MockLessonRepo) {
				repo.On("Update", 1, input).Return(1, nil).Once()
			},
			expectedID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := new(MockLessonRepo)
			tt.setup(lessonRepo)
			svc := NewLessonService(lessonRepo, new(MockCourseRepo), nil)

			id, err := svc.Update(1, input)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			lessonRepo.AssertExpectations(t)
		})
	}
}

func TestLessonService_DeleteByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		id         int
		setup      func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo)
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "lesson not found on get",
			id:   1,
			setup: func(lessonRepo *MockLessonRepo, _ *MockCourseRepo) {
				lessonRepo.On("GetByID", 1).Return(models.Lesson{}, models.ErrLessonNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "lesson not found",
			errWrapped: models.ErrLessonNotFound,
		},
		{
			name: "internal error on lesson get",
			id:   2,
			setup: func(lessonRepo *MockLessonRepo, _ *MockCourseRepo) {
				lessonRepo.On("GetByID", 2).Return(models.Lesson{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get lesson",
			errWrapped: errDB,
		},
		{
			name: "course not found",
			id:   3,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 3).Return(models.Lesson{ID: 3, CourseID: 100}, nil).Once()
				courseRepo.On("GetByID", ctx, 100).Return(models.Course{}, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal error on course get",
			id:   4,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 4).Return(models.Lesson{ID: 4, CourseID: 101}, nil).Once()
				courseRepo.On("GetByID", ctx, 101).Return(models.Course{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get course",
			errWrapped: errDB,
		},
		{
			name: "cannot delete from active course",
			id:   5,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 5).Return(models.Lesson{ID: 5, CourseID: 102}, nil).Once()
				courseRepo.On("GetByID", ctx, 102).Return(models.Course{ID: 102, IsActive: true}, nil).Once()
			},
			errCode: http.StatusConflict,
			errMsg:  "cannot delete lesson inside active course",
		},
		{
			name: "not found on delete",
			id:   6,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 6).Return(models.Lesson{ID: 6, CourseID: 103}, nil).Once()
				courseRepo.On("GetByID", ctx, 103).Return(models.Course{ID: 103, IsActive: false}, nil).Once()
				lessonRepo.On("DeleteByID", ctx, 6).Return(models.ErrLessonNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "lesson not found",
			errWrapped: models.ErrLessonNotFound,
		},
		{
			name: "internal error on delete",
			id:   7,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 7).Return(models.Lesson{ID: 7, CourseID: 104}, nil).Once()
				courseRepo.On("GetByID", ctx, 104).Return(models.Course{ID: 104, IsActive: false}, nil).Once()
				lessonRepo.On("DeleteByID", ctx, 7).Return(errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to delete lesson",
			errWrapped: errDB,
		},
		{
			name: "success",
			id:   8,
			setup: func(lessonRepo *MockLessonRepo, courseRepo *MockCourseRepo) {
				lessonRepo.On("GetByID", 8).Return(models.Lesson{ID: 8, CourseID: 105}, nil).Once()
				courseRepo.On("GetByID", ctx, 105).Return(models.Course{ID: 105, IsActive: false}, nil).Once()
				lessonRepo.On("DeleteByID", ctx, 8).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lessonRepo := new(MockLessonRepo)
			courseRepo := new(MockCourseRepo)
			tt.setup(lessonRepo, courseRepo)

			svc := NewLessonService(lessonRepo, courseRepo, nil)
			err := svc.DeleteByID(ctx, tt.id)

			if tt.errCode == 0 {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			lessonRepo.AssertExpectations(t)
			courseRepo.AssertExpectations(t)
		})
	}
}
