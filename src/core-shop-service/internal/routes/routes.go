package routes

import (
	"github.com/core/shop/golang/internal/repositories"
	"github.com/gin-gonic/gin"
)

func ProductManager(approachingRoute *gin.Engine) {
	approachingRoute.GET("products", repositories.GetProducts())
	approachingRoute.GET("products/product_id", repositories.GetProduct())
	approachingRoute.POST("products", repositories.InsertProducts())
	approachingRoute.PUT("products/product_id", repositories.UpdateProducts())
	approachingRoute.DELETE("products/product_id", repositories.DeleteProduct())
}
