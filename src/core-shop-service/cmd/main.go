package main

import (
	"fmt"
	"log"
	"os"

	"github.com/core/shop/golang/cmd/migrations"
	"github.com/core/shop/golang/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	if err := migrations.RunMigrations(); err != nil {
		fmt.Println(err)
	}
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

	routes.ProductManager(router)
	routes.OrdersManager(router)
	routes.UserManager(router)
	routes.CategoryManager(router)

	if err := router.Run(ipAddress + ":" + port); err != nil {
		fmt.Printf("Trying to run server on ip %s ...", ipAddress)
		log.Fatal("There's some error occured with running the server...")
	}

}
