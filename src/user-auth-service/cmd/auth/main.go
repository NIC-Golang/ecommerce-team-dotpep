package main

import (
	"fmt"
	"go/auth-service/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("/app/src/user-auth-service/.env")
	if err != nil {
		fmt.Print("error with loading .env file...")
	}
	router := gin.Default()
	port := os.Getenv("ENGINE_PORT")
	if port == "" {
		port = "8001"
	}

	ipAdress1, ipAdress2, ipAdress3 := os.Getenv("IP1"), os.Getenv("IP2"), os.Getenv("IP3")
	err = router.SetTrustedProxies([]string{ipAdress1, ipAdress2, ipAdress3})
	if err != nil {
		log.Fatal(err)
	}
	routes.AuthintificateRoute(router)
	routes.UserManager(router)
	if err := router.Run(ipAdress1 + ":" + port); err != nil {
		fmt.Printf("Trying to run the server on port %s..", ipAdress1)
		if err := router.Run(ipAdress2 + ":" + port); err != nil {
			fmt.Printf("Trying to run the server on port %s..", ipAdress2)
			if err := router.Run(ipAdress3 + ":" + port); err != nil {
				fmt.Printf("Trying to run the server on port %s..", ipAdress3)
				log.Fatal(err)
			}
		}
	}
}
