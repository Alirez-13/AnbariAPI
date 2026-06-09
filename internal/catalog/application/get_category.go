// # SINGLE REASON: Orchestrate category retrieval.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, id, true)
}
