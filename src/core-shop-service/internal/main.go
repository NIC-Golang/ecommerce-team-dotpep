package main

import (
	"os"

	"github.com/core/shop/golang/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.Default()
	router.Use(middleware.JWTAuthMiddleware())

	router.Run(":" + port)
}
