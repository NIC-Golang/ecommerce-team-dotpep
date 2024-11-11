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
	approachingRoute.DELETE("products/:product_id", repositories.DeleteProduc())
}

func OrdersManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("orders", repositories.GetOrders())
	approachingRoute.GET("orders/:order_id", repositories.GetOrder())
	approachingRoute.POST("orders", repositories.MakeAnOrder())
	approachingRoute.PUT("orders/:order_id", repositories.UpdateOrder())
	approachingRoute.DELETE("orders/:order_id", repositories.DeleteOrder())
}

func UserManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("users", repositories.GetUsers())
	approachingRoute.GET("users/:user_id", repositories.GetUser())
	approachingRoute.PUT("users/:user_id", repositories.UpdateUser())
	approachingRoute.DELETE("users/:user_id", repositories.DeleteUser())
}
