package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errDB                      = errors.New("db error")
	errBeginFailed             = errors.New("begin failed")
	errDeleteLessonsFailed     = errors.New("delete lessons failed")
	errDeleteEnrollmentsFailed = errors.New("delete enrollments failed")
	errCommitFailed            = errors.New("commit failed")
)

type mockCourseRepo struct{ mock.Mock }

type mockLessonRepo struct{ mock.Mock }

type mockEnrollmentRepo struct{ mock.Mock }

func (m *mockCourseRepo) GetAll() ([]models.Course, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.([]models.Course), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCourseRepo) GetByID(ctx context.Context, id int) (models.Course, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(models.Course), args.Error(1)
	}
	return models.Course{}, args.Error(1)
}

func (m *mockCourseRepo) DeleteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockCourseRepo) DeleteByIDTx(ctx context.Context, tx *sqlx.Tx, id int) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *mockCourseRepo) Create(input models.CreateCourse) (int, error) {
	args := m.Called(input)
	if v := args.Get(0); v != nil {
		return v.(int), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockCourseRepo) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {
	args := m.Called(ctx, id, input)
	if v := args.Get(0); v != nil {
		return v.(int), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockLessonRepo) GetAll() ([]models.Lesson, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.([]models.Lesson), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockLessonRepo) GetByID(id int) (models.Lesson, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(models.Lesson), args.Error(1)
	}
	return models.Lesson{}, args.Error(1)
}

func (m *mockLessonRepo) GetByCourseID(courseID int) ([]models.Lesson, error) {
	args := m.Called(courseID)
	if v := args.Get(0); v != nil {
		return v.([]models.Lesson), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockLessonRepo) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockLessonRepo) DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	args := m.Called(ctx, tx, courseID)
	return args.Error(0)
}

func (m *mockLessonRepo) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	args := m.Called(ctx, input)
	if v := args.Get(0); v != nil {
		return v.(int), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockLessonRepo) Update(id int, input models.UpdateLesson) (int, error) {
	args := m.Called(id, input)
	if v := args.Get(0); v != nil {
		return v.(int), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockEnrollmentRepo) Exists(ctx context.Context, userID, courseID int) (bool, error) {
	args := m.Called(ctx, userID, courseID)
	if v := args.Get(0); v != nil {
		return v.(bool), args.Error(1)
	}
	return false, args.Error(1)
}

func (m *mockEnrollmentRepo) Create(ctx context.Context, input models.CreateEnrollment) (int, error) {
	args := m.Called(ctx, input)
	if v := args.Get(0); v != nil {
		return v.(int), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockEnrollmentRepo) DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}

func (m *mockEnrollmentRepo) DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	args := m.Called(ctx, tx, courseID)
	return args.Error(0)
}

func (m *mockEnrollmentRepo) GetMyCourses(ctx context.Context, userID int) ([]models.MyCourse, error) {
	args := m.Called(ctx, userID)
	if v := args.Get(0); v != nil {
		return v.([]models.MyCourse), args.Error(1)
	}
	return nil, args.Error(1)
}

func ptr[T any](v T) *T { return &v }

func assertAppErr(t *testing.T, err error, expectedCode int, expectedMsg string, expectedWrapped error) {
	t.Helper()

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, expectedCode, appErr.Code)
	assert.Equal(t, expectedMsg, appErr.Message)
	if expectedWrapped != nil {
		assert.ErrorIs(t, err, expectedWrapped)
	}
}

func TestCourseService_Create(t *testing.T) {
	tests := []struct {
		name        string
		input       models.CreateCourse
		setup       func(*mockCourseRepo)
		expectedID  int
		errCode     int
		errMsg      string
		errWrapped  error
		errContains string
	}{
		{
			name:    "bad request validation",
			input:   models.CreateCourse{Title: "", Price: 100, TeacherID: 1},
			errCode: http.StatusBadRequest,
			errMsg:  "course title is required",
		},
		{
			name:  "teacher not found",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7},
			setup: func(r *mockCourseRepo) {
				r.On("Create", mock.MatchedBy(func(in models.CreateCourse) bool { return in.IsActive == false })).Return(0, models.ErrTeacherNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "teacher not found",
			errWrapped: models.ErrTeacherNotFound,
		},
		{
			name:  "slug conflict",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7},
			setup: func(r *mockCourseRepo) {
				r.On("Create", mock.Anything).Return(0, models.ErrSlugAlreadyExists).Once()
			},
			errCode:    http.StatusConflict,
			errMsg:     "slug already exists",
			errWrapped: models.ErrSlugAlreadyExists,
		},
		{
			name:  "internal error",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7},
			setup: func(r *mockCourseRepo) {
				r.On("Create", mock.Anything).Return(0, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to create course",
			errWrapped: errDB,
		},
		{
			name:  "success with forced inactive",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7, IsActive: true},
			setup: func(r *mockCourseRepo) {
				r.On("Create", mock.MatchedBy(func(in models.CreateCourse) bool {
					return in.Title == "Go" && in.IsActive == false
				})).Return(44, nil).Once()
			},
			expectedID: 44,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(mockCourseRepo)
			if tt.setup != nil {
				tt.setup(courseRepo)
			}
			svc := NewCourseService(courseRepo, new(mockLessonRepo), new(mockEnrollmentRepo), nil)

			id, err := svc.Create(tt.input)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			courseRepo.AssertExpectations(t)
		})
	}
}

func TestCourseService_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*mockCourseRepo)
		expected   []models.Course
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "success",
			setup: func(r *mockCourseRepo) {
				r.On("GetAll").Return([]models.Course{{ID: 1, Title: "Go"}}, nil).Once()
			},
			expected: []models.Course{{ID: 1, Title: "Go"}},
		},
		{
			name: "repo error",
			setup: func(r *mockCourseRepo) {
				r.On("GetAll").Return([]models.Course(nil), errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get courses",
			errWrapped: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(mockCourseRepo)
			tt.setup(courseRepo)
			svc := NewCourseService(courseRepo, new(mockLessonRepo), new(mockEnrollmentRepo), nil)

			courses, err := svc.GetAll()

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, courses)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			courseRepo.AssertExpectations(t)
		})
	}
}

func TestCourseService_GetByID(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		id         int
		setup      func(*mockCourseRepo)
		expected   models.Course
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "success",
			id:   1,
			setup: func(r *mockCourseRepo) {
				r.On("GetByID", ctx, 1).Return(models.Course{ID: 1, Title: "Go"}, nil).Once()
			},
			expected: models.Course{ID: 1, Title: "Go"},
		},
		{
			name: "not found",
			id:   2,
			setup: func(r *mockCourseRepo) {
				r.On("GetByID", ctx, 2).Return(models.Course{}, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal",
			id:   3,
			setup: func(r *mockCourseRepo) {
				r.On("GetByID", ctx, 3).Return(models.Course{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get course",
			errWrapped: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(mockCourseRepo)
			tt.setup(courseRepo)
			svc := NewCourseService(courseRepo, new(mockLessonRepo), new(mockEnrollmentRepo), nil)

			course, err := svc.GetByID(ctx, tt.id)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, course)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			courseRepo.AssertExpectations(t)
		})
	}
}

func TestCourseService_Update(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		id         int
		input      models.UpdateCourse
		setup      func(*mockCourseRepo, *mockLessonRepo)
		expectedID int
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name:  "activate check lessons error",
			id:    1,
			input: models.UpdateCourse{IsActive: ptr(true)},
			setup: func(_ *mockCourseRepo, l *mockLessonRepo) {
				l.On("GetByCourseID", 1).Return([]models.Lesson(nil), errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to check lessons",
			errWrapped: errDB,
		},
		{
			name:  "activate without lessons",
			id:    2,
			input: models.UpdateCourse{IsActive: ptr(true)},
			setup: func(_ *mockCourseRepo, l *mockLessonRepo) {
				l.On("GetByCourseID", 2).Return([]models.Lesson{}, nil).Once()
			},
			errCode:    http.StatusBadRequest,
			errMsg:     models.ErrCourseCannotBeActivated.Error(),
			errWrapped: models.ErrCourseCannotBeActivated,
		},
		{
			name:  "activate success",
			id:    3,
			input: models.UpdateCourse{IsActive: ptr(true)},
			setup: func(r *mockCourseRepo, l *mockLessonRepo) {
				l.On("GetByCourseID", 3).Return([]models.Lesson{{ID: 1}}, nil).Once()
				r.On("Update", ctx, 3, mock.MatchedBy(func(in models.UpdateCourse) bool {
					return in.IsActive != nil && *in.IsActive
				})).Return(3, nil).Once()
			},
			expectedID: 3,
		},
		{
			name:  "course not found",
			id:    4,
			input: models.UpdateCourse{Title: ptr("new")},
			setup: func(r *mockCourseRepo, _ *mockLessonRepo) {
				r.On("Update", ctx, 4, mock.MatchedBy(func(in models.UpdateCourse) bool {
					return in.Title != nil && *in.Title == "new"
				})).Return(0, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name:  "slug conflict",
			id:    5,
			input: models.UpdateCourse{Slug: ptr("dup")},
			setup: func(r *mockCourseRepo, _ *mockLessonRepo) {
				r.On("Update", ctx, 5, mock.Anything).Return(0, models.ErrSlugAlreadyExists).Once()
			},
			errCode:    http.StatusConflict,
			errMsg:     "slug already exists",
			errWrapped: models.ErrSlugAlreadyExists,
		},
		{
			name:  "internal",
			id:    6,
			input: models.UpdateCourse{Duration: ptr(10)},
			setup: func(r *mockCourseRepo, _ *mockLessonRepo) {
				r.On("Update", ctx, 6, mock.Anything).Return(0, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to update course",
			errWrapped: errDB,
		},
		{
			name:  "success no activation",
			id:    7,
			input: models.UpdateCourse{Duration: ptr(10)},
			setup: func(r *mockCourseRepo, _ *mockLessonRepo) {
				r.On("Update", ctx, 7, mock.Anything).Return(7, nil).Once()
			},
			expectedID: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(mockCourseRepo)
			lessonRepo := new(mockLessonRepo)
			if tt.setup != nil {
				tt.setup(courseRepo, lessonRepo)
			}
			svc := NewCourseService(courseRepo, lessonRepo, new(mockEnrollmentRepo), nil)

			id, err := svc.Update(ctx, tt.id, tt.input)

			if tt.errCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			courseRepo.AssertExpectations(t)
			lessonRepo.AssertExpectations(t)
		})
	}
}

func TestCourseService_DeleteByID(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		id         int
		setupDB    func(sqlmock.Sqlmock)
		setupMocks func(*mockCourseRepo, *mockLessonRepo, *mockEnrollmentRepo)
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "begin tx error",
			id:   1,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin().WillReturnError(errBeginFailed)
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to start transaction",
			errWrapped: errBeginFailed,
		},
		{
			name: "delete lessons error",
			id:   2,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectRollback()
			},
			setupMocks: func(_ *mockCourseRepo, l *mockLessonRepo, _ *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 2).Return(errDeleteLessonsFailed).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to delete related lessons",
			errWrapped: errDeleteLessonsFailed,
		},
		{
			name: "delete enrollments error",
			id:   3,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectRollback()
			},
			setupMocks: func(_ *mockCourseRepo, l *mockLessonRepo, e *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 3).Return(nil).Once()
				e.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 3).Return(errDeleteEnrollmentsFailed).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to delete related enrollments",
			errWrapped: errDeleteEnrollmentsFailed,
		},
		{
			name: "course not found",
			id:   4,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectRollback()
			},
			setupMocks: func(r *mockCourseRepo, l *mockLessonRepo, e *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 4).Return(nil).Once()
				e.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 4).Return(nil).Once()
				r.On("DeleteByIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 4).Return(models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "delete course internal",
			id:   5,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectRollback()
			},
			setupMocks: func(r *mockCourseRepo, l *mockLessonRepo, e *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 5).Return(nil).Once()
				e.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 5).Return(nil).Once()
				r.On("DeleteByIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 5).Return(errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to delete course",
			errWrapped: errDB,
		},
		{
			name: "commit error",
			id:   6,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectCommit().WillReturnError(errCommitFailed)
			},
			setupMocks: func(r *mockCourseRepo, l *mockLessonRepo, e *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 6).Return(nil).Once()
				e.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 6).Return(nil).Once()
				r.On("DeleteByIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 6).Return(nil).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to commit transaction",
			errWrapped: errCommitFailed,
		},
		{
			name: "success",
			id:   7,
			setupDB: func(sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectCommit()
			},
			setupMocks: func(r *mockCourseRepo, l *mockLessonRepo, e *mockEnrollmentRepo) {
				l.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 7).Return(nil).Once()
				e.On("DeleteByCourseIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 7).Return(nil).Once()
				r.On("DeleteByIDTx", ctx, mock.AnythingOfType("*sqlx.Tx"), 7).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, sm, err := sqlmock.New()
			require.NoError(t, err)
			defer sqlDB.Close()

			db := sqlx.NewDb(sqlDB, "sqlmock")
			courseRepo := new(mockCourseRepo)
			lessonRepo := new(mockLessonRepo)
			enrollRepo := new(mockEnrollmentRepo)
			if tt.setupDB != nil {
				tt.setupDB(sm)
			}
			if tt.setupMocks != nil {
				tt.setupMocks(courseRepo, lessonRepo, enrollRepo)
			}
			svc := NewCourseService(courseRepo, lessonRepo, enrollRepo, db)

			err = svc.DeleteByID(ctx, tt.id)

			if tt.errCode == 0 {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assertAppErr(t, err, tt.errCode, tt.errMsg, tt.errWrapped)
			}

			courseRepo.AssertExpectations(t)
			lessonRepo.AssertExpectations(t)
			enrollRepo.AssertExpectations(t)
			require.NoError(t, sm.ExpectationsWereMet())
		})
	}
}
