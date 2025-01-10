package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/notifier-service/internal/config"
	"github.com/notifier-service/internal/routes"
)

func main() {
	config.RabbitMQ()
	router := gin.Default()
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	router.SetTrustedProxies([]string{os.Getenv("IP")})
	routes.AuthRoutes(router)
	if err := router.Run(os.Getenv("IP") + ":" + os.Getenv("PORT")); err != nil {
		log.Printf("error running the server on ip: %s and port: %s", os.Getenv("IP"), os.Getenv("PORT"))
	}

}
