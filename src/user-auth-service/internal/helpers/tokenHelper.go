package helpers

import (
	"go/auth-service/internal/config"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email    string
	Name     string
	UserType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")
var key = os.Getenv("KEY")

func CreateToken(email, name, userType string) (tokenWithClaims, refreshTokenWithClaims string, err error) {
	claims := SignedDetails{
		Email:    email,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}

	refreshClaims := SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}

	tokenWithClaims, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshTokenWithClaims, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(key))
	if err != nil {
		log.Panic(err)
		return
	}

	return
}
