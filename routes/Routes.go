package routes

import (
	"AnbariAPI/Internal/Inventory"
	"AnbariAPI/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes registers all application routes.
// TODO (Caller): Update main.go to pass the *gorm.DB instance: routes.SetupRoutes(r, db)
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
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
		repo := Inventory.NewRepository(db)
		resolver := Inventory.NewExitResolver()
		invSvc := Inventory.NewInventoryService(repo, resolver)
		invH := Inventory.NewInventoryHandler(invSvc)

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
