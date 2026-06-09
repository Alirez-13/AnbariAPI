// # SINGLE REASON: Orchestrate category listing.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) ListCategories(ctx context.Context) ([]domain.Category, int64, error) {
	return s.repo.ListCategories(ctx, true)
}
