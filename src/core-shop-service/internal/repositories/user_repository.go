package repositories

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = 1$", adminID).Scan(&role)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "No client found with the specified ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
			}
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You have no rights for this action"})
			return
		}
		var user models.CLient
		err = conn.QueryRow(ctx, "SELECT client_id,client_name, client_last_name, client_email, client_phone, client_type,token, refresh_token FROM clients WHERE client_id = 1$", userId).Scan(
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

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		adminID := c.GetHeader("AdminID")
		err := godotenv.Load()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close(ctx)
		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = $1", adminID).Scan(&role)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusForbidden, gin.H{"error": "No client found with the specified Admin ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying admin role"})
			}
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You have no rights for this action!"})
			return
		}
		result, err := conn.Exec(ctx, "DELETE FROM clients WHERE client_id = $1", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No client found with the specified ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		adminID := c.GetHeader("AdminID")

		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		password, host := os.Getenv("SQL_PASS"), os.Getenv("HOST_SQL")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		connStr := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var role string
		err = conn.QueryRow(ctx, "SELECT client_type FROM clients WHERE client_id = $1", adminID).Scan(&role)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusForbidden, gin.H{"error": "No client found with the specified Admin ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying admin role"})
			}
			return
		}
		defer conn.Close(ctx)
		if role != "ADMIN" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You have no rights for this action!"})
			return
		}

		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		query := "UPDATE clients SET "
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
		query += fmt.Sprintf(" WHERE client_id = $%d", paramID)
		params = append(params, userId)

		result, err := conn.Exec(ctx, query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No client found with the specified ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}
