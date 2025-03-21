package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/events"
	"github.com/notifier-service/internal/models"
)

func AdminNotify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var admin models.Admin
		err := c.ShouldBindJSON(&admin)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if admin.Name == "" || admin.Email == "" || admin.Phone == "" {
			c.JSON(400, gin.H{"error": "name, email and phone are required"})
			return
		}

		events.AdminEvent(admin.Name, admin.Email, admin.Phone)
	}
}
