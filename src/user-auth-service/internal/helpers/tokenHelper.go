package helpers

import (
	"context"
	"fmt"
	"go/auth-service/internal/config"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Name     string
	UserType string
	Uid      string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var key = os.Getenv("KEY")

func CreateToken(email, name, userType, userId string) (tokenWithClaims, refreshTokenWithClaims string, err error) {
	claims := SignedDetails{
		Email:    email,
		Name:     name,
		UserType: userType,
		Uid:      userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshClaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
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

func UpdateTokens(token, refreshToken, userId string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	localzone := time.FixedZone("UTC+5", 5*60*60)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refreshToken})
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now().In(localzone).Format(time.RFC3339)})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		return fmt.Errorf("failed to update tokens for user %s: %v", userId, err)
	}

	return nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(
		signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "invalid token"
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, ""
}
