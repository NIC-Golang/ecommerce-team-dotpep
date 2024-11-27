package repositories

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/core/shop/golang/internal/config"
	"github.com/core/shop/golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
			return
		}
		defer conn.Close(ctx)

		rows, err := conn.Query(ctx, "SELECT client_id, client_name, client_last_name, client_email, client_phone, client_type, token, refresh_token FROM clients")
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

		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
			return
		}
		defer conn.Close(ctx)

		var user models.CLient
		err = conn.QueryRow(ctx, "SELECT client_id, client_name, client_last_name, client_email, client_phone, client_type, token, refresh_token FROM clients WHERE client_id = $1", id).Scan(
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
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		Id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close(ctx)

		result, err := conn.Exec(ctx, "DELETE FROM clients WHERE client_id = $1", Id)
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
		Id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		conn, err := config.GetDBConnection(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		defer conn.Close(ctx)

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
		params = append(params, Id)

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
