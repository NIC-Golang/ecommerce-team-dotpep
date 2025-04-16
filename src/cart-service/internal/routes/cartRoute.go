package routes

import (
	"cart-service/golang/internal/repositories"

	"github.com/gin-gonic/gin"
)

func CartManager(approachingRoute *gin.Engine) {

	cartRoutes := approachingRoute.Group("/cart")
	{
		cartRoutes.POST("", repositories.AddToCart())
		cartRoutes.GET("", repositories.GetCart())
		cartRoutes.DELETE("/:id", repositories.DeleteItemFromCart())
		cartRoutes.DELETE("", repositories.ClearCart())
		cartRoutes.GET("/:product_id", repositories.FindCartItemsByID())

	}
}

func OrderManager(route *gin.Engine) {
	orderRoute := route.Group("/order")
	{
		orderRoute.POST("/checkout", repositories.OrderCreating())
		orderRoute.GET("/:id", repositories.GetOrder())
		orderRoute.POST("/status:status", repositories.ChangeStatus())
	}
}
