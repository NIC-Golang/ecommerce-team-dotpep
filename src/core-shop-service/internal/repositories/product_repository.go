package repositories

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/core/shop/golang/internal/helpers"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := godotenv.Load(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !helpers.CheckUserType(c, "ADMIN") {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have the rights to perform this action"})
			return
		}

		host := os.Getenv("HOST_SQL")
		password := os.Getenv("SQL_PASS")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
			return
		}
		defer conn.Close(ctx)

		rows, err := conn.Query(ctx, "SELECT * FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The database query could not be executed"})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			err := rows.Scan(
				&product.ID,
				&product.Name,
				&product.Description,
				&product.Price,
				&product.SKU,
				&product.Quantity,
				&product.Created_at,
				&product.Update_at,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Product data scanning error"})
				return
			}
			products = append(products, product)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing product data"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
