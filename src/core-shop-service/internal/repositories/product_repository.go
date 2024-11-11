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
)

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetHeader("AdminID")
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
		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = $1", adminID).Scan(&role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have the rights to perform this action"})
			return
		}
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

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("product_id")
		adminID := c.GetHeader("AdminID")
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close(ctx)
		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = $1", adminID).Scan(&role)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid AdminID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching admin role"})
			}
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You have no rights for this action!"})
			return
		}
		var product models.Product
		err = conn.QueryRow(ctx, "SELECT product_id, product_name, product_description, product_price, product_sku, product_quantity, product_created_at, product_updated_at FROM products WHERE product_id = $1", productId).Scan(
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
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching product data"})
			}
			return
		}
		c.JSON(http.StatusOK, product)
	}
}
