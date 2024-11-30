package controllers

import (
	"context"
	"go/auth-service/internal/config"
	"go/auth-service/internal/helpers"
	"go/auth-service/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		userType := "USER"
		localzone := time.FixedZone("UTC+5", 5*60*60)
		token, refreshToken, err := helpers.CreateToken(*user.Email, *user.Name, *user.Type)
		if err != nil {
			log.Fatal(err)
			return
		}
		newUser := models.User{
			ID:           primitive.NewObjectID(),
			Name:         user.Name,
			Email:        user.Email,
			Phone:        user.Phone,
			Password:     user.Password,
			Type:         &userType,
			Token:        token,
			RefreshToken: refreshToken,
			Created_at:   time.Now().In(localzone),
			Updated_at:   time.Now().In(localzone),
		}

		resultInsertionNumber, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resultInsertionNumber)
	}
}
