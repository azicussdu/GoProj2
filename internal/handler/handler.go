package handler

import (
	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	courses := api.Group("/courses")
	{
		courses.GET("", h.GetCourses)
		courses.GET("/:id", h.GetCourseByID)
		courses.DELETE("/:id", h.DeleteCourse)
		courses.POST("", h.CreateCourse)
		courses.PUT("/:id", h.UpdateCourse)
	}

	lessons := api.Group("/lessons")
	{
		lessons.GET("", h.GetLessons)
		lessons.GET("/:id", h.GetLessonByID)
		lessons.DELETE("/:id", h.DeleteLesson)
		lessons.POST("", h.CreateLesson)
		lessons.PUT("/:id", h.UpdateLesson)
	}

	return r, nil
}
