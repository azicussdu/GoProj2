package handler

import (
	"errors"
	"net/http"

	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var input models.RegisterUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.services.Auth.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) Login(c *gin.Context) {
	var input models.LoginUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tokens, err := h.services.Auth.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {
	type refreshInput struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var input refreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tokens, err := h.services.Auth.Refresh(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) ChangeUserRole(c *gin.Context) {
	//userID, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
	//	return
	//}
	//
	//var input models.ChangeUserRoleInput
	//if err := c.ShouldBindJSON(&input); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	//	return
	//}
	//
	//updatedID, err := h.services.Auth.ChangeUserRole(userID, input.Role)
	//if err != nil {
	//	switch {
	//	case errors.Is(err, models.ErrUserNotFound):
	//		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	//		return
	//	case errors.Is(err, models.ErrInvalidRole):
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be teacher or admin"})
	//		return
	//	case errors.Is(err, models.ErrRoleChangeOnlyFromStudent):
	//		c.JSON(http.StatusConflict, gin.H{"error": "only student role can be changed"})
	//		return
	//	default:
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change user role"})
	//		return
	//	}
	//}

	//c.JSON(http.StatusOK, gin.H{"id": updatedID})
}
