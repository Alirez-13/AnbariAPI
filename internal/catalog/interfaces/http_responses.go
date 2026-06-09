// # SINGLE REASON: Map catalog domain entities to HTTP responses.
package interfaces

import (
	"AnbariAPI/internal/catalog/domain"
	"AnbariAPI/internal/catalog/dto"
)

func categoryResponse(category *domain.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func productResponse(product *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:         product.ID,
		CategoryID: product.CategoryID,
		Category:   categoryResponse(&product.Category),
		Name:       product.Name,
		Attribute:  product.Attribute,
		PackSize:   product.PackSize,
		CreatedAt:  product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
