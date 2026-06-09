// # SINGLE REASON: Preserve category route handler function names.
package handler

import (
	"AnbariAPI/internal/catalog/application"
	"AnbariAPI/internal/catalog/infrastructure"
	cataloghttp "AnbariAPI/internal/catalog/interfaces"
	"AnbariAPI/shared/database"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	defaultHandler().CreateCategory(c)
}

func GetCategory(c *gin.Context) {
	defaultHandler().GetCategory(c)
}

func ListCategories(c *gin.Context) {
	defaultHandler().ListCategories(c)
}

func DeleteCategory(c *gin.Context) {
	defaultHandler().DeleteCategory(c)
}

func defaultHandler() *cataloghttp.CatalogHandler {
	repo := infrastructure.NewGormRepository(database.GetDB())
	service := application.NewCatalogService(repo)
	return cataloghttp.NewCatalogHandler(service)
}
