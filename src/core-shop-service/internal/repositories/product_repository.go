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

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
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

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("product_id")
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

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("product_id")
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		connStr := fmt.Sprintf("postgres://Fivret:%s@localhost:%s/project", password, host)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result, err := conn.Exec(ctx, "DELETE FROM products WHERE product_id = $1", productId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No product found with the specified ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := c.Param("product_id")
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}
		defer conn.Close(ctx)

		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		if len(input) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided for update"})
			return
		}

		query := "UPDATE products SET "
		params := []interface{}{}
		paramID := 1

		for key, value := range input {
			if paramID > 1 {
				query += ", "
			}
			query += fmt.Sprintf("%s = $%d", key, paramID)
			params = append(params, value)
			paramID++
		}
		query += fmt.Sprintf(" WHERE product_id = $%d", paramID)
		params = append(params, productId)

		result, err := conn.Exec(ctx, query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No product found with the specified ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

func InsertProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		_ = godotenv.Load()
		password, port, host := os.Getenv("SQL_PASS"), os.Getenv("PORT_SQL"), os.Getenv("IP2")
		if port == "" {
			port = "5432"
		}
		connStr := fmt.Sprintf("postgres://Fiveret:%s@%s:%s/project?sslmode=disable", password, host, port)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}
		defer conn.Close(ctx)

		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		if len(input) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided for update"})
			return
		}

		query := "INSERT INTO products (product_name, product_description, product_price, product_sku, product_quantity, product_created_at, product_updated_at) "
		location := time.FixedZone("UTC+5", 5*60*60)
		created_at := time.Now().In(location).Format(time.RFC3339)
		updated_at := time.Now().In(location).Format(time.RFC3339)
		values := "VALUES ($1, $2, $3, $4, $5, $6, $7)"
		var params []interface{}
		params = append(params, input["product_name"], input["product_description"], input["product_price"], input["product_sku"], input["product_quantity"], created_at, updated_at)

		result, err := conn.Exec(ctx, query+values, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert product: %v", err)})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to insert product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product inserted successfully!"})

	}
}
