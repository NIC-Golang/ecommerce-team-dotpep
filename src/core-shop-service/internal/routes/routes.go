package routes

import (
	"github.com/core/shop/golang/internal/repositories"
	"github.com/gin-gonic/gin"
)

func ProductManager(approachingRoute *gin.Engine) {
	productRoutes := approachingRoute.Group("/products")
	{
		productRoutes.GET("", repositories.GetProducts())
		productRoutes.GET("/:product_id", repositories.GetProduct())
		productRoutes.POST("", repositories.InsertProduct())
		productRoutes.PUT("/:product_id", repositories.UpdateProduct())
		productRoutes.DELETE("/:product_id", repositories.DeleteProduct())
	}
}

func OrdersManager(approachingRoute *gin.Engine) {
	orderRoutes := approachingRoute.Group("/orders")
	{
		orderRoutes.GET("", repositories.GetOrders())
		orderRoutes.GET("/:client_id", repositories.GetUsersOrders())
		orderRoutes.POST("", repositories.MakeAnOrder())
		orderRoutes.DELETE("/:order_id", repositories.DeleteOrderByOrderId())
	}
}

func UserManager(approachingRoute *gin.Engine) {
	userRoutes := approachingRoute.Group("/users")
	{
		userRoutes.GET("", repositories.GetUsers())
		userRoutes.GET("/:user_id", repositories.GetUser())
		userRoutes.PUT("/:user_id", repositories.UpdateUser())
		userRoutes.DELETE("/:user_id", repositories.DeleteUser())
	}
}

func CategoryManager(approachingRoute *gin.Engine) {
	categoryRoutes := approachingRoute.Group("/categories")
	{
		categoryRoutes.GET("/:id", repositories.GetCategory())
		categoryRoutes.GET("", repositories.GetCategories())
		categoryRoutes.POST("", repositories.CreateCategory())
		categoryRoutes.DELETE("/:id", repositories.DeleteCategory())
	}
}
