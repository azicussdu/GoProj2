package middleware

import (
	"net/http"
	"strings"

	"github.com/azicussdu/GoProj2/internal/auth"
	"github.com/azicussdu/GoProj2/internal/models"
	"github.com/gin-gonic/gin"
)

const userContextKey = "auth_user"

func Auth(tmanager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		user, err := tmanager.ParseAccessToken(strings.TrimSpace(parts[1]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(userContextKey, *user)
		c.Next()
	}
}

func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {

	allowedMap := make(map[models.UserRole]struct{}, len(allowedRoles))

	for _, role := range allowedRoles {
		allowedMap[role] = struct{}{}
	}

	return func(c *gin.Context) {
		userVal, ok := c.Get(userContextKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
			return
		}

		user, ok := userVal.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user"})
			return
		}

		if _, ok := allowedMap[user.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role permissions"})
			return
		}

		c.Next()
	}
}
