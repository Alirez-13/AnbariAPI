package application

// # SINGLE REASON: Orchestrate inbound inventory transaction use case.

import (
	"context"

	models "AnbariAPI/internal/inventory/domain"
)

func (s *InventoryServiceImpl) ProcessInboundTransaction(ctx context.Context, input InboundDTO) (*models.Transaction, error) {
	if err := validateLines(input.Lines); err != nil {
		return nil, err
	}

	t := buildTransaction(TransactionTypeInbound, input.Reference, input.Date, input.Lines)
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.SaveTransaction(txCtx, t); err != nil {
			return err
		}

		batches := make([]*models.InventoryBatch, 0, len(t.Lines))
		for i := range t.Lines {
			line := &t.Lines[i]
			batches = append(batches, &models.InventoryBatch{
				TransactionLineID: line.ID,
				ProductID:         line.ProductID,
				InitialQuantity:   line.BaseQuantity,
				RemainingQuantity: line.BaseQuantity,
				EntryUnitCost:     line.UnitPrice.Div(line.UnitMultiplier),
				EntryDate:         t.Date,
				Version:           1,
			})
		}

		return s.repo.SaveInventoryBatches(txCtx, batches)
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}
