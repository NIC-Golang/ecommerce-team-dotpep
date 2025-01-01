package middleware

import (
	"go/auth-service/internal/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentification() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		clientToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}

func AdminRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if userType != "ADMIN" || !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to access this route!"})
			c.Abort()
			return
		}
	}
}
