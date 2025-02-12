package repositories

import (
	"cart-service/golang/internal/helpers"
	"cart-service/golang/internal/models"
	"cart-service/golang/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(500, gin.H{"error": "Header is missing!"})
			return
		}
		id, err := helpers.IdAuthorization(authHeader)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var items []models.CartItem
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		existingCart, err := redis.GetCartFromRedis(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if existingCart == nil {
			existingCart = &models.Cart{
				UserID:    id,
				CreatedAt: time.Now().In(time.FixedZone("UTC+5", 5*60*60)),
				Items:     []models.CartItem{},
			}
		}

		for _, newItem := range items {
			found := false

			for i, existingItem := range existingCart.Items {
				if existingItem.ProductID == newItem.ProductID {
					existingCart.Items[i].Quantity += newItem.Quantity
					found = true
					break
				}
			}
			if !found {
				existingCart.Items = append(existingCart.Items, newItem)
			}
		}

		existingCart.UpdatedAt = time.Now().In(time.FixedZone("UTC+5", 5*60*60))

		if err := redis.SaveToCart(id, existingCart); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save cart"})
			return
		}
		c.JSON(200, gin.H{"message": "Items added to cart", "cart": existingCart})
	}
}

func GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		cart, err := redis.GetCartFromRedis(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error retrieving cart", "details": err.Error()})
			return
		}
		if cart == nil {
			c.JSON(404, gin.H{"error": "Cart not found"})
			return
		}

		c.JSON(200, gin.H{"message": "Cart was found!", "cart": cart})
	}
}

func ClearCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		err = redis.DeleteCartFromRedis(id)
		if err != nil {
			if err.Error() == "cart not found" {
				c.JSON(404, gin.H{"error": "Cart not found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Your cart was deleted successfully!"})
	}
}

// func DeleteItemFromCart() gin.HandlerFunc{
// 	return func(c *gin.Context){
// 		productId := c.Param("id")

// 	}
// }
