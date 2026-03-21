package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCourseRepo struct {
	mock.Mock
}

type MockLessonRepo struct {
	mock.Mock
}

type MockEnrollmentRepo struct {
	mock.Mock
}

var _ repository.CourseRepo = (*MockCourseRepo)(nil)
var _ repository.LessonRepo = (*MockLessonRepo)(nil)
var _ repository.EnrollmentRepo = (*MockEnrollmentRepo)(nil)

func (m *MockCourseRepo) GetAll() ([]models.Course, error) {
	args := m.Called()

	var courses []models.Course
	if value := args.Get(0); value != nil {
		courses = value.([]models.Course)
	}

	return courses, args.Error(1)
}

func (m *MockCourseRepo) GetByID(ctx context.Context, id int) (models.Course, error) {
	args := m.Called(ctx, id)

	var course models.Course
	if value := args.Get(0); value != nil {
		course = value.(models.Course)
	}

	return course, args.Error(1)
}

func (m *MockCourseRepo) DeleteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCourseRepo) DeleteByIDTx(ctx context.Context, tx *sqlx.Tx, id int) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockCourseRepo) Create(input models.CreateCourse) (int, error) {
	args := m.Called(input)
	return args.Int(0), args.Error(1)
}

func (m *MockCourseRepo) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {
	args := m.Called(ctx, id, input)
	return args.Int(0), args.Error(1)
}

func (m *MockLessonRepo) GetAll() ([]models.Lesson, error) {
	args := m.Called()

	var lessons []models.Lesson
	if value := args.Get(0); value != nil {
		lessons = value.([]models.Lesson)
	}

	return lessons, args.Error(1)
}

func (m *MockLessonRepo) GetByID(id int) (models.Lesson, error) {
	args := m.Called(id)

	var lesson models.Lesson
	if value := args.Get(0); value != nil {
		lesson = value.(models.Lesson)
	}

	return lesson, args.Error(1)
}

func (m *MockLessonRepo) GetByCourseID(courseID int) ([]models.Lesson, error) {
	args := m.Called(courseID)

	var lessons []models.Lesson
	if value := args.Get(0); value != nil {
		lessons = value.([]models.Lesson)
	}

	return lessons, args.Error(1)
}

func (m *MockLessonRepo) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLessonRepo) DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	args := m.Called(ctx, tx, courseID)
	return args.Error(0)
}

func (m *MockLessonRepo) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), args.Error(1)
}

func (m *MockLessonRepo) Update(id int, input models.UpdateLesson) (int, error) {
	args := m.Called(id, input)
	return args.Int(0), args.Error(1)
}

func (m *MockEnrollmentRepo) Exists(ctx context.Context, userID, courseID int) (bool, error) {
	args := m.Called(ctx, userID, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEnrollmentRepo) Create(ctx context.Context, input models.CreateEnrollment) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), args.Error(1)
}

func (m *MockEnrollmentRepo) DeleteByUserAndCourse(ctx context.Context, userID, courseID int) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}

func (m *MockEnrollmentRepo) DeleteByCourseIDTx(ctx context.Context, tx *sqlx.Tx, courseID int) error {
	args := m.Called(ctx, tx, courseID)
	return args.Error(0)
}

func (m *MockEnrollmentRepo) GetMyCourses(ctx context.Context, userID int) ([]models.MyCourse, error) {
	args := m.Called(ctx, userID)

	var courses []models.MyCourse
	if value := args.Get(0); value != nil {
		courses = value.([]models.MyCourse)
	}

	return courses, args.Error(1)
}

//---------------------------------------------------------------------

func assertAppErr(t *testing.T, err error, expectedCode int, expectedMessage string, expectedWrapped error) {
	t.Helper()

	var appErr *apperror.AppError

	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, expectedCode, appErr.Code)
	assert.Equal(t, expectedMessage, appErr.Message)

	if expectedWrapped != nil {
		assert.ErrorIs(t, err, expectedWrapped)
	}
}

var errDB = errors.New("db error")

//---------------------------------------------------------------------

func TestCourseService_GetByID_Success(t *testing.T) {
	mCourseRepo := &MockCourseRepo{} // = new(MockCourseRepo)
	mLessonRepo := &MockLessonRepo{}
	mEnrollRepo := &MockEnrollmentRepo{}

	ctx := context.Background()
	expectedCourse := models.Course{ID: 1, Title: "Golang"}

	mCourseRepo.On("GetByID", ctx, 1).Return(expectedCourse, nil).Once()

	courseSrv := NewCourseService(mCourseRepo, mLessonRepo, mEnrollRepo, nil)
	course, err := courseSrv.GetByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, expectedCourse, course)

	mCourseRepo.AssertExpectations(t)
}

func TestCourseService_GetByID_NotFound(t *testing.T) {
	mCourseRepo := new(MockCourseRepo)
	mLessonRepo := new(MockLessonRepo)
	mEnrollRepo := new(MockEnrollmentRepo)

	ctx := context.Background()

	mCourseRepo.On("GetByID", ctx, 2).
		Return(models.Course{}, models.ErrCourseNotFound).Once()

	courseSrv := NewCourseService(mCourseRepo, mLessonRepo, mEnrollRepo, nil)
	_, err := courseSrv.GetByID(ctx, 2)

	require.Error(t, err)

	assertAppErr(
		t,
		err,
		http.StatusNotFound,
		"course not found",
		models.ErrCourseNotFound,
	)

	mCourseRepo.AssertExpectations(t)
}

func TestCourseService_GetByID_InternalError(t *testing.T) {
	mCourseRepo := new(MockCourseRepo)
	mLessonRepo := new(MockLessonRepo)
	mEnrollRepo := new(MockEnrollmentRepo)

	ctx := context.Background()

	mCourseRepo.On("GetByID", ctx, 3).
		Return(models.Course{}, errDB).Once()

	courseSrv := NewCourseService(mCourseRepo, mLessonRepo, mEnrollRepo, nil)
	_, err := courseSrv.GetByID(ctx, 3)

	require.Error(t, err)

	assertAppErr(
		t,
		err,
		http.StatusInternalServerError,
		"failed to get course",
		errDB,
	)

	mCourseRepo.AssertExpectations(t)
}

func TestSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "test 1",
			input:    2,
			expected: 4,
		},
		{
			name:     "test 2",
			input:    3,
			expected: 9,
		},
		{
			name:     "test 3",
			input:    5,
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input * tt.input
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TABLE-DRIVEN tests
func TestCourseService_GetByID(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		id         int
		setup      func(repo *MockCourseRepo)
		expected   models.Course
		errCode    int
		errMsg     string
		errWrapped error
	}{
		{
			name: "success",
			id:   1,
			setup: func(r *MockCourseRepo) {
				r.On("GetByID", ctx, 1).Return(models.Course{ID: 1, Title: "Go"}, nil).Once()
			},
			expected: models.Course{ID: 1, Title: "Go"},
		},
		{
			name: "not found",
			id:   2,
			setup: func(r *MockCourseRepo) {
				r.On("GetByID", ctx, 2).Return(models.Course{}, models.ErrCourseNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "course not found",
			errWrapped: models.ErrCourseNotFound,
		},
		{
			name: "internal",
			id:   3,
			setup: func(r *MockCourseRepo) {
				r.On("GetByID", ctx, 3).Return(models.Course{}, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to get course",
			errWrapped: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(MockCourseRepo)
			tt.setup(courseRepo)
			svc := NewCourseService(courseRepo, new(MockLessonRepo), new(MockEnrollmentRepo), nil)

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

func TestCourseService_Create(t *testing.T) {
	tests := []struct {
		name        string
		input       models.CreateCourse
		setup       func(repo *MockCourseRepo)
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
			setup: func(r *MockCourseRepo) {
				r.On("Create",
					mock.MatchedBy(func(in models.CreateCourse) bool { return in.IsActive == false })).
					Return(0, models.ErrTeacherNotFound).Once()
			},
			errCode:    http.StatusNotFound,
			errMsg:     "teacher not found",
			errWrapped: models.ErrTeacherNotFound,
		},
		{
			name:  "slug conflict",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7},
			setup: func(r *MockCourseRepo) {
				r.On("Create", mock.Anything).Return(0, models.ErrSlugAlreadyExists).Once()
			},
			errCode:    http.StatusConflict,
			errMsg:     "slug already exists",
			errWrapped: models.ErrSlugAlreadyExists,
		},
		{
			name:  "internal error",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7},
			setup: func(r *MockCourseRepo) {
				r.On("Create", mock.Anything).Return(0, errDB).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMsg:     "failed to create course",
			errWrapped: errDB,
		},
		{
			name:  "success with forced inactive",
			input: models.CreateCourse{Title: "Go", Slug: "go", Price: 100, TeacherID: 7, IsActive: true},
			setup: func(r *MockCourseRepo) {
				r.On("Create", mock.MatchedBy(func(in models.CreateCourse) bool {
					return in.Title == "Go" && in.IsActive == false
				})).Return(44, nil).Once()
			},
			expectedID: 44,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseRepo := new(MockCourseRepo)
			if tt.setup != nil {
				tt.setup(courseRepo)
			}
			svc := NewCourseService(courseRepo, new(MockLessonRepo), new(MockEnrollmentRepo), nil)

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
