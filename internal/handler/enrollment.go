package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) EnrollCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	userVal, ok := c.Get("auth_user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user"})
		return
	}

	enrollmentID, err := h.services.Enrollment.JoinCourse(c.Request.Context(), user, courseID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrOnlyStudentsCanEnroll):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		case errors.Is(err, models.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		case errors.Is(err, models.ErrEnrollmentAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enroll in course"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": enrollmentID})
}

func (h *Handler) LeaveCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	userVal, ok := c.Get("auth_user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user"})
		return
	}

	err = h.services.Enrollment.LeaveCourse(c.Request.Context(), user, courseID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrOnlyStudentsCanEnroll):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		case errors.Is(err, models.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		case errors.Is(err, models.ErrEnrollmentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to leave course"})
			return
		}
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "left course"})
}

func (h *Handler) GetMyCourses(c *gin.Context) {
	userVal, ok := c.Get("auth_user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user"})
		return
	}

	coursesEnrolled, err := h.services.Enrollment.GetMyCourses(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, models.ErrOnlyStudentsCanEnroll) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get my courses"})
		return
	}

	c.JSON(http.StatusOK, coursesEnrolled)
}
