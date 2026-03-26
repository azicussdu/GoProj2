package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCourseService struct {
	mock.Mock
}

var _ service.Course = (*MockCourseService)(nil)

func (m *MockCourseService) Create(input models.CreateCourse) (int, error) {
	args := m.Called(input)
	return args.Int(0), args.Error(1)
}

func (m *MockCourseService) GetAll() ([]models.Course, error) {
	args := m.Called()

	var courses []models.Course
	if v := args.Get(0); v != nil {
		courses = v.([]models.Course)
	}

	return courses, args.Error(1)
}

func (m *MockCourseService) GetByID(ctx context.Context, id int) (models.Course, error) {
	args := m.Called(ctx, id)

	var course models.Course
	if v := args.Get(0); v != nil {
		course = v.(models.Course)
	}

	return course, args.Error(1)
}

func (m *MockCourseService) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseService) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {
	args := m.Called(ctx, id, input)
	return args.Int(0), args.Error(1)
}

// -----------------------------------

func ptr[T any](v T) *T {
	return &v
}

func newCourseHandler(courseSvc service.Course) *Handler {
	return &Handler{
		services: &service.Services{
			Course: courseSvc,
		},
	}
}

func newTestContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)                                                  // отключаем лишние логи (чтобы не мешать логам тестов)
	w := httptest.NewRecorder()                                                // фейковый HTTP-ответ
	c, _ := gin.CreateTestContext(w)                                           // создаем фейк контекст и привязывем к - w (handler → пишет в c → это попадает в w)
	c.Request = httptest.NewRequest(method, path, bytes.NewBufferString(body)) // как будто пришёл реальный HTTP запрос
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/json") // c.ShouldBindJSON(...) чтобы работал нормально
	}
	return c, w
}

func decodeBodyMap(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func TestHandler_GetCourses(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(svc *MockCourseService)
		errCode    int
		errMessage string
	}{
		{
			name: "success",
			setup: func(svc *MockCourseService) {
				svc.On("GetAll").Return([]models.Course{{ID: 1, Title: "Go"}}, nil).Once()
			},
		},
		{
			name: "service error",
			setup: func(svc *MockCourseService) {
				svc.On("GetAll").Return(nil, apperror.Internal("failed to get courses", assert.AnError)).Once()
			},
			errCode:    http.StatusInternalServerError,
			errMessage: "failed to get courses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(MockCourseService)
			tt.setup(svc)
			h := newCourseHandler(svc)

			c, w := newTestContext(http.MethodGet, "/courses", "")
			h.GetCourses(c)

			if tt.errCode == 0 {
				require.Equal(t, http.StatusOK, w.Code)
				var got []models.Course
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, []models.Course{{ID: 1, Title: "Go"}}, got)
			} else {
				require.Equal(t, tt.errCode, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, tt.errMessage, resp["error"])
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_GetCourseByID(t *testing.T) {
	tests := []struct {
		name       string
		idParam    string
		setup      func(svc *MockCourseService)
		errCode    int
		errMessage string
	}{
		{
			name:       "invalid id",
			idParam:    "abc",
			errCode:    http.StatusBadRequest,
			errMessage: "invalid course id",
		},
		{
			name:    "not found",
			idParam: "2",
			setup: func(svc *MockCourseService) {
				svc.On("GetByID", mock.Anything, 2).Return(models.Course{}, apperror.NotFound("course not found", models.ErrCourseNotFound)).Once()
			},
			errCode:    http.StatusNotFound,
			errMessage: "course not found",
		},
		{
			name:    "success",
			idParam: "1",
			setup: func(svc *MockCourseService) {
				svc.On("GetByID", mock.Anything, 1).Return(models.Course{ID: 1, Title: "Go"}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(MockCourseService)
			if tt.setup != nil {
				tt.setup(svc)
			}
			h := newCourseHandler(svc)

			c, w := newTestContext(http.MethodGet, "/courses/"+tt.idParam, "")
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}
			h.GetCourseByID(c)

			if tt.errCode == 0 {
				require.Equal(t, http.StatusOK, w.Code)
				var got models.Course
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, models.Course{ID: 1, Title: "Go"}, got)
			} else {
				require.Equal(t, tt.errCode, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, tt.errMessage, resp["error"])
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateCourse(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setup      func(svc *MockCourseService)
		errCode    int
		errMessage string
	}{
		{
			name:       "invalid body",
			body:       `{"title":`,
			errCode:    http.StatusBadRequest,
			errMessage: "invalid request body",
		},
		{
			name: "service conflict",
			body: `{"title":"Go","slug":"go","price":10,"teacher_id":1}`,
			setup: func(svc *MockCourseService) {
				svc.On("Create", models.CreateCourse{Title: "Go", Slug: "go", Price: 10, TeacherID: 1}).Return(0, apperror.Conflict("slug already exists", models.ErrSlugAlreadyExists)).Once()
			},
			errCode:    http.StatusConflict,
			errMessage: "slug already exists",
		},
		{
			name: "success",
			body: `{"title":"Go","slug":"go","price":10,"teacher_id":1}`,
			setup: func(svc *MockCourseService) {
				svc.On("Create", models.CreateCourse{Title: "Go", Slug: "go", Price: 10, TeacherID: 1}).Return(11, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(MockCourseService)
			if tt.setup != nil {
				tt.setup(svc)
			}
			h := newCourseHandler(svc)

			c, w := newTestContext(http.MethodPost, "/courses", tt.body)
			h.CreateCourse(c)

			if tt.errCode == 0 {
				require.Equal(t, http.StatusCreated, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, float64(11), resp["id"])
			} else {
				require.Equal(t, tt.errCode, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, tt.errMessage, resp["error"])
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdateCourse(t *testing.T) {
	tests := []struct {
		name       string
		idParam    string
		body       string
		setup      func(svc *MockCourseService)
		errCode    int
		errMessage string
	}{
		{
			name:       "invalid id",
			idParam:    "abc",
			body:       `{"title":"Go"}`,
			errCode:    http.StatusBadRequest,
			errMessage: "invalid course id",
		},
		{
			name:       "invalid body",
			idParam:    "1",
			body:       `{"title":`,
			errCode:    http.StatusBadRequest,
			errMessage: "invalid request body",
		},
		{
			name:    "service error",
			idParam: "1",
			body:    `{"is_active":true}`,
			setup: func(svc *MockCourseService) {
				svc.On("Update", mock.Anything, 1, models.UpdateCourse{IsActive: ptr(true)}).Return(0, apperror.BadRequest(models.ErrCourseCannotBeActivated.Error(), models.ErrCourseCannotBeActivated)).Once()
			},
			errCode:    http.StatusBadRequest,
			errMessage: models.ErrCourseCannotBeActivated.Error(),
		},
		{
			name:    "success",
			idParam: "1",
			body:    `{"title":"Go v2"}`,
			setup: func(svc *MockCourseService) {
				svc.On("Update", mock.Anything, 1, models.UpdateCourse{Title: ptr("Go v2")}).Return(1, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(MockCourseService)
			if tt.setup != nil {
				tt.setup(svc)
			}
			h := newCourseHandler(svc)

			c, w := newTestContext(http.MethodPut, "/courses/"+tt.idParam, tt.body)
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}
			h.UpdateCourse(c)

			if tt.errCode == 0 {
				require.Equal(t, http.StatusOK, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, float64(1), resp["id"])
			} else {
				require.Equal(t, tt.errCode, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, tt.errMessage, resp["error"])
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteCourse(t *testing.T) {
	tests := []struct {
		name       string
		idParam    string
		setup      func(svc *MockCourseService)
		errCode    int
		errMessage string
	}{
		{
			name:       "invalid id",
			idParam:    "abc",
			errCode:    http.StatusBadRequest,
			errMessage: "invalid course id",
		},
		{
			name:    "service not found",
			idParam: "2",
			setup: func(svc *MockCourseService) {
				svc.On("DeleteByID", mock.Anything, 2).Return(apperror.NotFound("course not found", models.ErrCourseNotFound)).Once()
			},
			errCode:    http.StatusNotFound,
			errMessage: "course not found",
		},
		{
			name:    "success",
			idParam: "1",
			setup: func(svc *MockCourseService) {
				svc.On("DeleteByID", mock.Anything, 1).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(MockCourseService)
			if tt.setup != nil {
				tt.setup(svc)
			}
			h := newCourseHandler(svc)

			c, w := newTestContext(http.MethodDelete, "/courses/"+tt.idParam, "")
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}
			h.DeleteCourse(c)

			if tt.errCode == 0 {
				require.Equal(t, http.StatusNoContent, w.Code)
				assert.Equal(t, "", w.Body.String())
			} else {
				require.Equal(t, tt.errCode, w.Code)
				resp := decodeBodyMap(t, w)
				assert.Equal(t, tt.errMessage, resp["error"])
			}

			svc.AssertExpectations(t)
		})
	}
}
