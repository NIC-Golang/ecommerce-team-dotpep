package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/events"
)

type order struct {
	Description string `json:"description"`
	Email       string `json:"email"`
}

func OrderCreating() gin.HandlerFunc {
	return func(c *gin.Context) {
		name, exists1 := c.Get("name")
		email, exists2 := c.Get("email")
		if !exists1 || !exists2 {
			c.JSON(500, gin.H{"error": "Name or email not found in context"})
			return
		}

		var order order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("%s, your order was created successfully!", name)})

		events.OrderEvent(name.(string), order.Description, email.(string))
	}
}
