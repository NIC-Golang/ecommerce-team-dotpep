package routes

import (
	"github.com/core/shop/golang/internal/repositories"
	"github.com/gin-gonic/gin"
)

func ProductManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("products", repositories.GetProducts())
	approachingRoute.GET("products/:product_id", repositories.GetProduct())
	approachingRoute.POST("products", repositories.InsertProduct())
	approachingRoute.PUT("products/:product_id", repositories.UpdateProduct())
	approachingRoute.DELETE("products/:product_id", repositories.DeleteProduct())
}

func OrdersManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("orders", repositories.GetOrders())
	approachingRoute.GET("orders/:client_id", repositories.GetUsersOrders())
	approachingRoute.POST("orders", repositories.MakeAnOrder())
	approachingRoute.DELETE("orders/:order_id", repositories.DeleteOrderByOrderId())
}

func UserManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("users", repositories.GetUsers())
	approachingRoute.GET("users/:user_id", repositories.GetUser())
	approachingRoute.PUT("users/:user_id", repositories.UpdateUser())
	approachingRoute.DELETE("users/:user_id", repositories.DeleteUser())
}
