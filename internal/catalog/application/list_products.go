// # SINGLE REASON: Orchestrate product listing.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) ListProducts(ctx context.Context) ([]domain.Product, int64, error) {
	return s.repo.ListProducts(ctx, true)
}
