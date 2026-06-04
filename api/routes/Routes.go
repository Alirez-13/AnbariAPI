package routes

import (
	handler2 "AnbariAPI/Internal/Inventory/handler"
	"AnbariAPI/Internal/Inventory/repository"
	"AnbariAPI/Internal/Inventory/resolver"
	"AnbariAPI/Internal/Inventory/service"
	handler3 "AnbariAPI/Internal/auth/handler"
	"AnbariAPI/Internal/catalog/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes registers all application routes.
func SetupRoutes(r *gin.Engine, db *gorm.DB, authHandler *handler3.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}

	v1 := r.Group("/api/v1")
	{
		categories := v1.Group("/categories")
		{
			// Note: Consider injecting handlers here in the future instead of using package-level functions.
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

		// Dependency Injection: Wire up the Inventory bounded context
		repo := repository.NewRepository(db)
		resolver := resolver.NewExitResolver()
		invSvc := service.NewInventoryService(repo, resolver)
		invH := handler2.NewInventoryHandler(invSvc)

		inventory := v1.Group("")
		{
			// Batch availability (for exit popup)
			inventory.GET("/products/:productId/batches", invH.GetAvailableBatches)

			// Transactions
			inventory.POST("/transactions/entry", invH.ProcessEntry)
			inventory.POST("/transactions/exit/preview", invH.PreviewExit)
			inventory.POST("/transactions/exit", invH.ConfirmExit)
		}
	}
}
