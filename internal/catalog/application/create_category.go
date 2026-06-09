// # SINGLE REASON: Orchestrate category creation.
package application

import (
	"context"
	"fmt"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) CreateCategory(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}

	category := &domain.Category{Name: input.Name}
	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
