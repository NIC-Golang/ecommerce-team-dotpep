package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/core/shop/golang/internal/config"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
)

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		}

		rows, err := conn.Query(ctx, "SELECT * from categories")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The database query could not be executed"})
			return
		}
		var categories []models.Category
		for rows.Next() {
			var category models.Category
			err := rows.Scan(
				&category.ID,
				&category.Name,
				&category.CreatedAt,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Product data scanning error"})
				return
			}

			categories = append(categories, category)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing product data"})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot connect to database..."})
			return
		}
		defer conn.Close(ctx)

		var input struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		localzone := time.FixedZone("UTC+5", 5*60*60)
		createdAt := time.Now().In(localzone).Format(time.RFC3339)

		query := "INSERT INTO categories (name, created_at) VALUES ($1, $2)"
		params := []interface{}{input.Name, createdAt}

		result, err := conn.Exec(ctx, query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert category: %v", err)})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to insert category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category inserted successfully!"})
	}
}

func DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot connect to database..."})
			return
		}
		defer conn.Close(ctx)
		result, err := conn.Exec(ctx, "DELETE FROM categories WHERE id = $1", categoryId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert category: %v", err)})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to insert category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully!"})
	}
}
