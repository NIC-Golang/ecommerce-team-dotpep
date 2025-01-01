package middleware

import (
	"go/auth-service/internal/helpers"
	"log"
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
		userToken, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			return
		}

		tokenStr, ok := userToken.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		log.Printf("Received token: %s", tokenStr)

		claims, err := helpers.ExtractClaimsFromToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if claims.UserType != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this route!"})
			return
		}
	}
}
