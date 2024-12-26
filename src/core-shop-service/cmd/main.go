package main

import (
	"fmt"
	"log"
	"os"

	//"github.com/core/shop/golang/internal/middleware"
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
	ipAddress1 := os.Getenv("IP1")
	ipAddress2 := os.Getenv("IP2")
	ipAddress3 := os.Getenv("IP3")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{ipAddress2, ipAddress1, ipAddress, ipAddress3})
	//router.Use(middleware.JWTAuthMiddleware())

	routes.ProductManager(router)
	routes.OrdersManager(router)
	routes.UserManager(router)
	routes.CategoryManager(router)

	if err := router.Run(ipAddress + ":" + port); err != nil {
		fmt.Printf("Trying to run server on ip %s ...", ipAddress)
		if err := router.Run(ipAddress1 + ":" + port); err != nil {
			fmt.Printf("Trying to run server on ip %s ...", ipAddress1)
			if err := router.Run(ipAddress2 + ":" + port); err != nil {
				fmt.Printf("Trying to run server on ip %s ...", ipAddress2)
				if err := router.Run(ipAddress3 + ":" + port); err != nil {
					fmt.Printf("Trying to run server on ip %s ...", ipAddress3)
					log.Fatal("There's some error occured with running the server...")
				}
			}
		}
	}
}
