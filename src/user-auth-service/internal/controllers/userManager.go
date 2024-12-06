package controllers

import (
	"context"
	"go/auth-service/internal/config"
	"go/auth-service/internal/helpers"
	"go/auth-service/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollect *mongo.Collection = config.GetCollection(config.DB, "users")

const timeoutDuration = 5 * time.Second

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		var user models.User
		err := userCollect.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		msg := helpers.CheckType(c, userId)
		if msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
