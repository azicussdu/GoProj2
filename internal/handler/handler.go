package handler

import (
	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService
	lessonService *service.LessonService
	// userService *service.UserService
	//TODO add other service

}

func NewHandler(cs *service.CourseService, ls *service.LessonService) *Handler {
	return &Handler{
		courseService: cs,
		lessonService: ls,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()

	r.GET("/courses", h.GetCourses)
	r.GET("/courses/:id", h.GetCourseByID) // localhost:8080/courses/@#@
	r.DELETE("/courses/:id", h.DeleteCourse)
	r.POST("/courses", h.CreateCourse)
	r.PUT("/courses/:id", h.UpdateCourse)

	r.GET("/lessons", h.GetLessons)
	r.GET("/lessons/:id", h.GetLessonByID)
	r.DELETE("/lessons/:id", h.DeleteLesson)
	r.POST("/lessons", h.CreateLesson)
	r.PUT("/lessons/:id", h.UpdateLesson)

	return r, nil
}
