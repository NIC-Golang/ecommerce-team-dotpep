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
	err := godotenv.Load()
	if err != nil {
		fmt.Print("error with loading .env file...")
	}
	router := gin.Default()
	port := os.Getenv("ENGINE_PORT")
	if port == "" {
		port = "8000"
	}

	ipAdress1, apAdress2, ipAdress3 := os.Getenv("IP1"), os.Getenv("IP2"), os.Getenv("IP3")
	err = router.SetTrustedProxies([]string{ipAdress1, apAdress2, ipAdress3})
	if err != nil {
		log.Fatal(err)
	}
	routes.AuthintificateRoute(router)
	routes.UserManager(router)
	if err := router.Run(ipAdress1 + ":" + port); err != nil {
		fmt.Print("Trying to run the server...")
		log.Fatal(err)
	}
}
