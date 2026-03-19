package service

import (
	"context"
	"testing"

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

func TestCourseService_GetByID_Success(t *testing.T) {
	mCourseRepo := new(MockCourseRepo)
	mLessonRepo := new(MockLessonRepo)
	mEnrollRepo := new(MockEnrollmentRepo)

	ctx := context.Background()
	expectedCourse := models.Course{ID: 1, Title: "Golang"}

	mCourseRepo.On("GetByID", ctx, 1).Return(expectedCourse, nil)

	courseSrv := NewCourseService(mCourseRepo, mLessonRepo, mEnrollRepo, nil)
	course, err := courseSrv.GetByID(ctx, 2)

	require.NoError(t, err)
	assert.Equal(t, expectedCourse, course)

	mCourseRepo.AssertExpectations(t)
}
