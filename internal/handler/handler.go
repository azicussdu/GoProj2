package handler

import (
	"github.com/azicussdu/GoProj2/internal/auth"
	"github.com/azicussdu/GoProj2/internal/middleware"
	"github.com/azicussdu/GoProj2/internal/models"
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

	api := r.Group("/api")

	// AUTH ROUTES
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
	}

	// PUBLIC COURSES
	courses := api.Group("/courses")
	{
		courses.GET("", h.GetCourses)
		courses.GET("/:id", h.GetCourseByID)
	}

	// PROTECTED ROUTES
	protected := api.Group("/")
	protected.Use(middleware.Auth(h.tokenManager))
	{
		courses := protected.Group("/courses")
		courses.Use(middleware.RequireRole(models.RoleTeacher, models.RoleAdmin))
		{
			courses.POST("", h.CreateCourse)
			courses.PUT("/:id", h.UpdateCourse)
			courses.DELETE("/:id", h.DeleteCourse)
		}

		// protected lesson actions
		lessons := protected.Group("/lessons")
		{
			lessons.POST("", middleware.RequireRole(models.RoleTeacher, models.RoleAdmin), h.CreateLesson)
			lessons.PUT("/:id", h.UpdateLesson)
			lessons.DELETE("/:id", middleware.RequireRole(models.RoleTeacher, models.RoleAdmin), h.DeleteLesson)
		}

		users := protected.Group("/users")
		{
			users.PUT("/:id/role", middleware.RequireRole(models.RoleAdmin), h.ChangeUserRole)
		}
	}

	return r, nil
}
