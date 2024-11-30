package helpers

import (
	"go/auth-service/internal/config"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email    string
	Name     string
	UserType string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")
var key = os.Getenv("KEY")

func CreateToken(email, name, userType string) (tokenWithClaims, refreshTokenWithClaims string, err error) {
	claims := SignedDetails{
		Email:    email,
		Name:     name,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Установка времени истечения токена
		},
	}

	refreshClaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)), // Установка времени истечения refresh токена
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenWithClaims, err = token.SignedString([]byte(key))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenWithClaims, err = refreshToken.SignedString([]byte(key))
	if err != nil {
		log.Panic(err)
		return
	}

	return
}
