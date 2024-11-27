package routes

import (
	"github.com/gin-gonic/gin"
)

func UserManager(aproachingRoute *gin.Engine) {
	userRoutes := aproachingRoute.Group("/users")
	{
		userRoutes.Use(middleware.Authentification())
		userRoutes.GET("", repositories.GetAll())
		userRoutes.GET("/:id", repositories.GetUser())
	}
}
