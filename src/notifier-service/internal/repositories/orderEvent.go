package repositories

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type order struct {
	Id string `json:"order_id"`
}

func OrderCreating() gin.HandlerFunc {
	return func(c *gin.Context) {
		name, exists := c.Get("name")
		if !exists {
			c.JSON(500, gin.H{"error": "Name not found in context"})
			return
		}
		var order order
		err := c.ShouldBindJSON(&order)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("%s , your order created successfully, here's your order id: %s", name, order.Id)})
	}
}
