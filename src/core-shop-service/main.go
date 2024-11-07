package main

import (
	"log"
	"net/http"
	"os"

	//"github.com/core/shop/golang/internal/middleware"
	"github.com/core/shop/golang/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ipAddress := os.Getenv("IP")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{ipAddress})
	//router.Use(middleware.JWTAuthMiddleware())

	routes.ProductManager(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "You successfully on api-1!"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "You successfully on api-2!"})
	})
	if err := router.Run(ipAddress + ":" + port); err != nil {
		log.Fatal("There's some error occured with running the server...")
	}
}
