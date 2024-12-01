package routes

import (
	controllers "go/auth-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

func AuthintificateRoute(approachingRoute *gin.Engine) {
	authRoutes := approachingRoute.Group("/users")
	{
		authRoutes.POST("/login", controllers.Login())
		authRoutes.POST("/signup", controllers.SignUp())
	}
}
