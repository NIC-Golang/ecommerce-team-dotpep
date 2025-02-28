package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/events"
)

type order struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}

func OrderCreating() gin.HandlerFunc {
	return func(c *gin.Context) {

		var order order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("%s, your order was created successfully!", order.Name)})

		events.OrderEvent(order.Name, order.Description, order.Email)
	}
}
