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
	//type refreshInput struct {
	//	RefreshToken string `json:"refresh_token" binding:"required"`
	//}
	//
	//var input refreshInput
	//if err := c.ShouldBindJSON(&input); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	//	return
	//}
	//
	//tokens, err := h.services.Auth.Refresh(input.RefreshToken)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, tokens)
}
