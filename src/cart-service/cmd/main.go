package main

import (
	"cart-service/golang/internal/redis"
	"cart-service/golang/internal/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("error downloading env file")
	}
	ipAdress, addrRedis := os.Getenv("IP"), os.Getenv("adrRedis")
	err = redis.InitRedis(addrRedis)
	if err == nil {
		fmt.Println("Connected to Redis")
	}
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("CART_SERVICE_PORT")
	if port == "" {
		port = "8003"
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{ipAdress})
	routes.CartManager(router)
	err = router.Run(ipAdress + ":" + port)
	if err != nil {
		fmt.Printf("Trying to run server on ip %s ...", ipAdress)
		log.Fatal("There's some error occured with running the server...")
	}
}
