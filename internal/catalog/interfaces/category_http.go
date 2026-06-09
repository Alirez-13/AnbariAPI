// # SINGLE REASON: Handle category HTTP requests.
package interfaces

import (
	"net/http"
	"strconv"

	"AnbariAPI/internal/catalog/application"
	"AnbariAPI/internal/catalog/dto"

	"github.com/gin-gonic/gin"
)

func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	category, err := h.service.CreateCategory(c.Request.Context(), application.CreateCategoryInput{Name: req.Name})
	if err != nil {
		writeCatalogError(c, err, "Failed to create category")
		return
	}
	c.JSON(http.StatusCreated, categoryResponse(category))
}

func (h *CatalogHandler) GetCategory(c *gin.Context) {
	id, ok := parseID(c, "Invalid category ID")
	if !ok {
		return
	}
	category, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		writeCatalogError(c, err, "Category not found")
		return
	}
	products := make([]dto.ProductResponse, len(category.Products))
	for i := range category.Products {
		products[i] = productResponse(&category.Products[i])
	}
	c.JSON(http.StatusOK, gin.H{"id": category.ID, "name": category.Name, "products": products, "created_at": category.CreatedAt, "updated_at": category.UpdatedAt})
}

func (h *CatalogHandler) ListCategories(c *gin.Context) {
	categories, total, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err, "Failed to fetch categories")
		return
	}
	responses := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = categoryResponse(&categories[i])
	}
	c.JSON(http.StatusOK, dto.CategoryListResponse{Categories: responses, Total: total})
}

func (h *CatalogHandler) DeleteCategory(c *gin.Context) {
	id, ok := parseID(c, "Invalid category ID")
	if !ok {
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		writeCatalogError(c, err, "Category not found")
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Category deleted successfully"})
}

func parseID(c *gin.Context, message string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: message})
		return 0, false
	}
	return uint(id), true
}
