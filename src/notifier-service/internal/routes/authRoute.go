package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/notifier-service/internal/repositories"
)

func AuthRoutes(incomingRoute *gin.Engine) {
	authEventRoute := incomingRoute.Group("/auth")
	{
		authEventRoute.POST("/signup", repositories.SignUp())
		authEventRoute.POST("/login", repositories.Login())
	}

}

func OrderRoutes(incomingRoute *gin.Engine) {
	orderEventRoute := incomingRoute.Group("/orders")
	{
		orderEventRoute.POST("", repositories.OrderCreating())
	}
}
