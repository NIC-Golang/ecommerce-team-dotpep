package main

import (
	"fmt"
	"log"
	"os"

	"github.com/core/shop/golang/internal/middleware"
	"github.com/core/shop/golang/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load("/app/.env")
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
	router.Use(middleware.AdminAuth())

	routes.ProductManager(router)
	routes.OrdersManager(router)
	routes.UserManager(router)
	routes.CategoryManager(router)

	if err := router.Run(ipAddress + ":" + port); err != nil {
		fmt.Printf("Trying to run server on ip %s ...", ipAddress)
		log.Fatal("There's some error occured with running the server...")
	}

}
