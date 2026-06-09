// # SINGLE REASON: Orchestrate product retrieval.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	return s.repo.GetProductByID(ctx, id, true)
}
