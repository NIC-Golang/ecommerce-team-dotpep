package routes

import (
	"cart-service/golang/internal/repositories"

	"github.com/gin-gonic/gin"
)

func CartManager(approachingRoute *gin.Engine) {
	cartRoutes := approachingRoute.Group("/cart")
	{
		cartRoutes.POST("/orders", repositories.AddToCart())
		cartRoutes.GET("", repositories.GetCart())
		cartRoutes.DELETE(":id", repositories.DeleteItemFromCart())
		cartRoutes.DELETE("", repositories.ClearCart())
		cartRoutes.GET(":product_id", repositories.FindCartItemsByID())
	}
}
