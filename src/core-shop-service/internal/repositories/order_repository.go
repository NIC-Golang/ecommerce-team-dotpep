package repositories

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := godotenv.Load(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading .env file"})
			return
		}
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
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
				&order.CreatedAt,
				&order.UpdatedAt,
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
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
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
				&order.CreatedAt,
				&order.UpdatedAt,
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
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
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
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There was some error with connection to database..."})
			return
		}
		defer conn.Close(ctx)

		tx, err := conn.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
			return
		}
		defer tx.Rollback(ctx)

		query := "INSERT INTO orders (user_id, total_amount, status, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING order_id"
		err = tx.QueryRow(ctx, query, order.UserID, order.TotalAmount, order.Status).Scan(&order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting order"})
			return
		}

		for _, item := range order.Products {
			_, err := tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
				order.ID, item.ProductID, item.Quantity, item.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting order item"})
				return
			}
		}

		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order_id": order.ID})
	}
}
