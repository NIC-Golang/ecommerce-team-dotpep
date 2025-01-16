package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/events"
	"github.com/notifier-service/internal/models"
)

func OrderCreating() gin.HandlerFunc {
	return func(c *gin.Context) {
		name, exists1 := c.Get("name")
		email, exists2 := c.Get("email")
		if !exists1 || !exists2 {
			c.JSON(500, gin.H{"error": "Name not found in context"})
			return
		}
		var order models.Order
		err := c.ShouldBindJSON(&order)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("%s , your order created successfully, here's your order id: %d", name, order.ID)})
		events.OrderEvent(name.(string), fmt.Sprintf("%d", order.ID), email.(string))
	}
}
