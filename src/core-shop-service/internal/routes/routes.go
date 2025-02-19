package routes

import (
	"github.com/core/shop/golang/internal/middleware"
	"github.com/core/shop/golang/internal/repositories"
	"github.com/gin-gonic/gin"
)

func ProductManager(approachingRoute *gin.Engine) {
	productRoutes := approachingRoute.Group("/products")
	{
		productRoutes.Use(middleware.AdminAuth())
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
		orderRoutes.POST("", repositories.MakeAnOrder())
	}
}

func CategoryManager(approachingRoute *gin.Engine) {
	categoryRoutes := approachingRoute.Group("/categories")
	{
		categoryRoutes.Use(middleware.AdminAuth())
		categoryRoutes.GET("/:id", repositories.GetCategory())
		categoryRoutes.GET("", repositories.GetCategories())
		categoryRoutes.POST("", repositories.CreateCategory())
		categoryRoutes.DELETE("/:id", repositories.DeleteCategory())
	}
}
