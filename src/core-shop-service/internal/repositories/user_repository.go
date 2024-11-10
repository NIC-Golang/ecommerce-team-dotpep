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
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetHeader("AdminID")
		if adminID == "" {
			c.JSON(500, gin.H{"error": "Your header is empty!"})
			c.Abort()
			return
		}

		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		host := os.Getenv("HOST_SQL")
		password := os.Getenv("SQL_PASS")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		connect := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connect)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to connect to the database"})
			return
		}
		defer conn.Close(ctx)

		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = $1", adminID).Scan(&role)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "No client found with the specified ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
			}
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You have no rights for that"})
			return
		}

		rows, err := conn.Query(ctx, "SELECT * FROM clients")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []models.CLient
		for rows.Next() {
			var user models.CLient
			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.LastName,
				&user.Email,
				&user.Phone,
				&user.Type,
				&user.Token,
				&user.RefreshToken,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		adminID := c.GetHeader("AdminID")
		if err := godotenv.Load(); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with connecting to database..."})
			return
		}
		defer conn.Close(ctx)

		if !helpers.CheckUserType(c, "ADMIN") {
			c.JSON(http.StatusForbidden, gin.H{"error": "You have no rights for this action"})
			return
		}

	}
}
