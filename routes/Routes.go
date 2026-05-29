package routes

import (
	"AnbariAPI/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.POST("", handler.CreateCategory)
			categories.GET("", handler.ListCategories)
			categories.GET("/:id", handler.GetCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
		}

		products := v1.Group("/products")
		{
			products.POST("", handler.CreateProduct)
			products.GET("", handler.ListProducts)
			products.GET("/:id", handler.GetProduct)
			products.DELETE("/:id", handler.DeleteProduct)
		}
	}
}
