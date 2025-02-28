package repositories

import (
	"cart-service/golang/internal/helpers"
	"cart-service/golang/internal/models"
	"cart-service/golang/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
)

var localzone = time.FixedZone("UTC+5", 5*60*60)

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
				CreatedAt: time.Now().In(localzone),
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

		existingCart.UpdatedAt = time.Now().In(localzone)
		existingCart.UserID = id
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
func FindCartItemsByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		productId := c.Param("product_id")
		if productId == "" {
			c.JSON(500, gin.H{"error": "Hard to find product id"})
			return
		}

		item, err := redis.FindCartItem(productId, id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}

		c.JSON(200, item)
	}
}

func DeleteItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.IdAuthorization(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		productId := c.Param("id")

		cart, err := redis.GetCartFromRedis(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch cart"})
			return
		}

		var newItems []models.CartItem
		for _, item := range cart.Items {
			if item.ProductID != productId {
				newItems = append(newItems, item)
			}
		}

		if len(newItems) == len(cart.Items) {
			c.JSON(404, gin.H{"error": "Item not found in cart"})
			return
		}

		cart.Items = newItems
		cart.UpdatedAt = time.Now().In(localzone)
		err = redis.SaveToCart(id, cart)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update cart"})
			return
		}

		c.JSON(200, gin.H{"message": "Item deleted successfully"})
	}
}
