package handler

import (
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
		respondError(c, err)
		return
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
		respondError(c, err)
		return
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

	courses, err := h.services.Enrollment.GetMyCourses(c.Request.Context(), user)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, courses)
}
