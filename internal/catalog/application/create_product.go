// # SINGLE REASON: Orchestrate product creation.
package application

import (
	"context"
	"fmt"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	if input.CategoryID == 0 || input.Name == "" {
		return nil, fmt.Errorf("%w: category_id and name are required", ErrValidation)
	}
	if _, err := s.repo.GetCategoryByID(ctx, input.CategoryID, false); err != nil {
		return nil, err
	}

	product := &domain.Product{
		CategoryID: input.CategoryID,
		Name:       input.Name,
		Attribute:  input.Attribute,
		PackSize:   input.PackSize,
	}
	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	return s.repo.GetProductByID(ctx, product.ID, true)
}
