package handler

import (
	"net/http"

	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService
	//TODO add other service
}

func NewHandler(cs *service.CourseService) *Handler {
	return &Handler{
		courseService: cs,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()

	r.GET("/courses", h.GetCourses)
	// kop marshrut bolady

	return r, nil
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.courseService.GetAll()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, courses)
}
