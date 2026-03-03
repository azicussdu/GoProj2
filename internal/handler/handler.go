package handler

import (
	"github.com/azicussdu/GoProj2/internal/auth"
	"github.com/azicussdu/GoProj2/internal/middleware"
	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()
	//r.Use(gin.Logger())

	api := r.Group("/api")

	auth := api.Group("/auth") // api/auth/login
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
	}

	courses := api.Group("/courses")
	//courses.Use(gin.Logger())
	courses.Use(middleware.Auth(h.tokenManager))
	{
		courses.GET("", h.GetCourses)
		courses.GET("/:id", h.GetCourseByID) // GET api/courses/3
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
