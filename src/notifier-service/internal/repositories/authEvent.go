package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/events"
	"github.com/notifier-service/internal/models"
)

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if user.Name == "" || user.Email == "" {
			c.JSON(400, gin.H{"error": "name and email are required"})
			return
		}
		events.SignUpEvent(user.Name, user.Email)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUser models.User
		err := c.ShouldBindJSON(&loginUser)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if loginUser.Name == "" {
			c.JSON(400, gin.H{"error": "name is required"})
			return
		}
		events.LoginEvent(loginUser.Name, loginUser.Email)
	}
}
