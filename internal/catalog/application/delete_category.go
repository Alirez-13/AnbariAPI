// # SINGLE REASON: Orchestrate category deletion.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

func (s *CatalogService) DeleteCategory(ctx context.Context, id uint) error {
	deleted, err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return domain.ErrNotFound
	}
	return nil
}
