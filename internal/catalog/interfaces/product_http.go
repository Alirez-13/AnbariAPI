// # SINGLE REASON: Handle product HTTP requests.
package interfaces

import (
	"errors"
	"net/http"

	"AnbariAPI/internal/catalog/application"
	"AnbariAPI/internal/catalog/domain"
	"AnbariAPI/internal/catalog/dto"

	"github.com/gin-gonic/gin"
)

func (h *CatalogHandler) CreateProduct(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	product, err := h.service.CreateProduct(c.Request.Context(), application.CreateProductInput{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Attribute:  req.Attribute,
		PackSize:   req.PackSize,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: "Category not found"})
			return
		}
		writeCatalogError(c, err, "Category not found")
		return
	}
	c.JSON(http.StatusCreated, productResponse(product))
}

func (h *CatalogHandler) GetProduct(c *gin.Context) {
	id, ok := parseID(c, "Invalid product ID")
	if !ok {
		return
	}
	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		writeCatalogError(c, err, "Product not found")
		return
	}
	c.JSON(http.StatusOK, productResponse(product))
}

func (h *CatalogHandler) ListProducts(c *gin.Context) {
	products, total, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err, "Failed to fetch products")
		return
	}
	responses := make([]dto.ProductResponse, len(products))
	for i := range products {
		responses[i] = productResponse(&products[i])
	}
	c.JSON(http.StatusOK, dto.ProductListResponse{Products: responses, Total: total})
}

func (h *CatalogHandler) DeleteProduct(c *gin.Context) {
	id, ok := parseID(c, "Invalid product ID")
	if !ok {
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		writeCatalogError(c, err, "Product not found")
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Product deleted successfully"})
}
