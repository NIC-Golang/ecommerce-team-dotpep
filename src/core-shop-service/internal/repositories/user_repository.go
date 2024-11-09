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

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		host := os.Getenv("HOST_SQL")
		password := os.Getenv("SQL_PASS")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		connect := fmt.Sprintf("postgres://Fiveret:%s@localhost:%s/project", password, host)

		conn, err := pgx.Connect(ctx, connect)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was some error with connection to database..."})
			return
		}
		defer conn.Close(ctx)
		if !helpers.CheckUserType(c, "ADMIN") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You have no rights for that..."})
			return
		}

		rows, err := conn.Query(ctx, "SELECT * from clients")
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
				&user.Last_name,
				&user.Email,
				&user.Phone,
				&user.Type,
				&user.Token,
				&user.Refresh_token,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan products data"})
			return
		}
		c.JSON(http.StatusOK, users)

	}
}
