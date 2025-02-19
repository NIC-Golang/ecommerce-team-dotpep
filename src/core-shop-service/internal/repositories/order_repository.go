package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/core/shop/golang/internal/helpers"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
)

func MakeAnOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}
		token, err := helpers.HeaderTrimming(authHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		_, email, err := helpers.GetIdAndEmailFromToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		name, err := helpers.GetName(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var order []models.OrderItem
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		orderJSON, err := json.Marshal(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling order to JSON"})
			return
		}

		resp1, err := helpers.SendWithHeaders(authHeader, orderJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to cart-service"})
			return
		}
		defer resp1.Body.Close()

		if resp1.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp1.Body)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong with cart-service", "details": string(body)})
			return
		}
		var descriptionBuilder strings.Builder
		for _, item := range order {
			descriptionBuilder.WriteString(fmt.Sprintf("ProductID: %s, Quantity: %d, Price: %.2f\n", item.ProductID, item.Quantity, item.Price))
		}

		body := models.Notifier{
			Name:        name,
			Email:       email,
			Description: descriptionBuilder.String(),
		}
		bodyNotify, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling notifier body to JSON"})
			return
		}

		resp, err := helpers.SendNotifierRequest("http://notifier-service:8082/orders", http.MethodPost, bodyNotify)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to notifier-service"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong with notifier-service", "details": string(body)})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
	}
}
