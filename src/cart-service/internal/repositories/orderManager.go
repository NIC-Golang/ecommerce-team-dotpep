package repositories

import (
	"cart-service/golang/internal/helpers"
	"cart-service/golang/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func OrderCreating() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		cart, err := redis.GetCartFromRedis(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve cart"})
			return
		}
		err = redis.CreateOrder(id, cart, time.Now().In(localzone))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Order created successfully!"})
	}
}

type user struct {
	UserID string `json:"id"`
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id user
		if err := c.ShouldBindJSON(&id); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		order, err := redis.GetOrderFromRedis(id.UserID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Can not find order"})
			return
		}

		c.JSON(201, order)
	}
}

func ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		status := c.Param("status")
		if status == "" {
			c.JSON(500, gin.H{"error": "empty status parameter"})
			return
		}
		order, err := redis.GetOrderFromRedis(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		order.Status = status

		err = redis.SaveToOrder(id, order)
		if err != nil {
			c.JSON(500, gin.H{"error": "error during saving the order"})
			return
		}
		c.JSON(201, gin.H{"message": "Order's status updated successfully!"})
	}
}
