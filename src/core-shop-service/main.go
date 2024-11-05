package main

import (
	"log"
	"net/http"
	"os"

	//"github.com/core/shop/golang/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.Default()
	//router.Use(middleware.JWTAuthMiddleware())
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "You successfully on api-1!"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "You successfully on api-2!"})
	})
	if err := router.Run(":" + port); err != nil {
		log.Fatal("There's some error occured with running the server...")
	}
}
