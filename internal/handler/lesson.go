package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLessons(c *gin.Context) {
	lessons, err := h.services.Lesson.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select lessons"})
		return
	}

	c.JSON(http.StatusOK, lessons)
}

func (h *Handler) GetLessonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	lesson, err := h.services.Lesson.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *Handler) CreateLesson(c *gin.Context) {
	var input models.CreateLesson
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.services.Lesson.Create(input)
	if err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lesson"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var input models.UpdateLesson
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedID, err := h.services.Lesson.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson to update not found"})
			return
		case errors.Is(err, models.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedID})
}

func (h *Handler) DeleteLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	err = h.services.Lesson.DeleteByID(id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson to delete not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "lesson is deleted"})
}
