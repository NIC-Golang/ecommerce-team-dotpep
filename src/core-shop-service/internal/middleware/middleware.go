package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(500, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}
		resp, err := http.Post("http://user-auth-service:8081/validate-token", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to send request to user-auth-service"})
			return
		}

		if resp != nil && resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		role, ok := result["user_type"].(string)
		if !ok || role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access restricted to administrators"})
			c.Abort()
			return
		}

		c.Next()
	}
}
