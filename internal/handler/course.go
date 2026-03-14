package handler

import (
	"net/http"
	"strconv"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	var input models.UpdateCourse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatedID, err := h.services.Course.Update(c.Request.Context(), id, input)
	if err != nil {
		//if errors.Is(err, models.ErrCourseNotFound) {
		//	c.JSON(http.StatusNotFound, gin.H{"error": "course to update not found"})
		//	return
		//}
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedID})
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var input models.CreateCourse

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.services.Course.Create(input)
	if err != nil {
		//var status int
		//var message string
		//
		//switch {
		//case errors.Is(err, models.ErrTeacherNotFound):
		//	status = http.StatusNotFound
		//	message = err.Error()
		//case errors.Is(err, models.ErrSlugAlreadyExists):
		//	status = http.StatusConflict
		//	message = err.Error()
		//default:
		//	status = http.StatusInternalServerError
		//	message = "failed to create a lesson"
		//}
		//
		//c.JSON(status, gin.H{"error": message})
		//return

		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	err = h.services.Course.DeleteByID(c.Request.Context(), id)
	if err != nil {
		//if errors.Is(err, models.ErrCourseNotFound) {
		//	c.JSON(http.StatusNotFound, gin.H{"error": "course to delete not found"})
		//	return
		//}
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return

		respondError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "course is deleted"})
}

func (h *Handler) GetCourseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	course, err := h.services.Course.GetByID(c.Request.Context(), id)
	if err != nil {
		//if errors.Is(err, models.ErrCourseNotFound) {
		//	c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		//	return
		//}
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return

		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.services.Course.GetAll()
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, courses)
}
