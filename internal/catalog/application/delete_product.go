// # SINGLE REASON: Orchestrate product deletion.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) DeleteProduct(ctx context.Context, id uint) error {
	deleted, err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return domain.ErrNotFound
	}
	return nil
}
