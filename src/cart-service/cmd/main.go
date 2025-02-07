package main

import (
	"cart-service/golang/internal/config"
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
	_, err = config.RedisConnection(addrRedis)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{ipAdress})

	err = router.Run(ipAdress +":" + port)
	if err != nil {
		fmt.Printf("Trying to run server on ip %s ...", ipAdress)
		log.Fatal("There's some error occured with running the server...")
	}
}
