package repositories

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/core/shop/golang/internal/helpers"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
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

		if !helpers.CheckUserType(c, "ADMIN") {
			c.JSON(400, gin.H{"error": "You have no rights for this action"})
			return
		}
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

			itemsRows, err := conn.Query(ctx, "SELECT * FROM order_items WHERE order_id = $1", order.ID)
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
