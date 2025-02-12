package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/core/shop/golang/internal/config"
	"github.com/core/shop/golang/internal/helpers"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There was some error with connection to database..."})
			return
		}
		defer conn.Close(ctx)

		rows, err := conn.Query(ctx, "SELECT * FROM orders")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
			return
		}

		var orders []models.Order
		defer rows.Close()

		for rows.Next() {
			var order models.Order
			err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.TotalAmount,
				&order.Status,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			itemsRows, err := conn.Query(ctx, "SELECT product_id, quantity,price FROM order_items WHERE order_id = $1", order.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var items []models.OrderItem
			for itemsRows.Next() {
				var item models.OrderItem
				err := itemsRows.Scan(
					&item.ProductID,
					&item.Quantity,
					&item.Price,
				)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				items = append(items, item)
			}
			order.Products = items
			orders = append(orders, order)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan orders"})
			return
		}

		c.JSON(http.StatusOK, models.OrdersResponse{Orders: orders})
	}
}

func GetUsersOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("client_id")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There was some error with connection to database..."})
			return
		}
		defer conn.Close(ctx)

		rows, err := conn.Query(ctx, "SELECT order_id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE user_id = $1", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var orders []models.Order
		for rows.Next() {
			var order models.Order
			err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.TotalAmount,
				&order.Status,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			productRows, err := conn.Query(ctx, "SELECT product_id, quantity, price FROM order_items WHERE order_id = $1", order.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer productRows.Close()

			var products []models.OrderItem
			for productRows.Next() {
				var item models.OrderItem
				err := productRows.Scan(
					&item.ProductID,
					&item.Quantity,
					&item.Price,
				)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				products = append(products, item)
			}
			order.Products = products
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func DeleteOrderByOrderId() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There was some error with connection to database..."})
			return
		}
		defer conn.Close(ctx)

		result, err := conn.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order items"})
			return
		}
		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No items found for the specified order ID"})
			return
		}

		result, err = conn.Exec(ctx, "DELETE FROM orders WHERE order_id = $1", orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order"})
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No order found with the specified order ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
}

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
