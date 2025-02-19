package main

import (
	"fmt"
	"log"
	"os"

	"go/auth-service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	router := gin.Default()

	port := os.Getenv("ENGINE_PORT")
	if port == "" {
		port = "8001"
	}

	ips := []string{os.Getenv("IP1"), os.Getenv("IP2"), os.Getenv("IP3")}
	var trustedIPs []string
	for _, ip := range ips {
		if ip != "" {
			trustedIPs = append(trustedIPs, ip)
		}
	}
	if len(trustedIPs) == 0 {
		log.Fatal("No trusted IP addresses provided")
	}

	if err := router.SetTrustedProxies(trustedIPs); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	routes.AuthintificateRoute(router)
	routes.UserManager(router)

	var runErr error
	for _, ip := range trustedIPs {
		addr := fmt.Sprintf("%s:%s", ip, port)
		log.Printf("Trying to run server on %s...", addr)
		runErr = router.Run(addr)
		if runErr == nil {
			break
		}
		log.Printf("Failed to run on %s: %v", addr, runErr)
	}
	if runErr != nil {
		log.Fatal("Server failed to start on any provided IP:", runErr)
	}
}
