package handler

import (
	"errors"
	"net/http"

	"github.com/azicussdu/GoProj2/internal/apperror"
	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	var appErr *apperror.AppError

	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
	})
}
