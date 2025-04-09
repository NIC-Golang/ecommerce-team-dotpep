package controllers

import (
	"context"
	"go/auth-service/internal/helpers"
	"go/auth-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		authHeader := c.Request.Header.Get("Authorization")

		err := helpers.CheckAdmin(authHeader)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
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

func GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		err := helpers.CheckAdmin(authHeader)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()
		var users []models.User
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching users from the database"})
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user models.User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user data"})
				return
			}
			users = append(users, user)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

type notifyUser struct {
	NotifierId string `json:"user_id"`
}

func GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		if email == "" {
			c.JSON(500, gin.H{"error": "empty email provided"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()
		var foundUser notifyUser
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			c.JSON(500, gin.H{"error": "Couldn't find user with email"})
			return
		}
		c.JSON(200, foundUser)
	}
}
