package application

// # SINGLE REASON: Orchestrate outbound inventory transaction use case.

import (
	"context"
	"errors"
	"fmt"

	models "AnbariAPI/internal/inventory/domain"

	"github.com/shopspring/decimal"
)

func (s *InventoryServiceImpl) ProcessOutboundTransaction(ctx context.Context, input OutboundDTO) (*models.Transaction, error) {
	if err := validateLines(input.Lines); err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 1; attempt <= s.retries; attempt++ {
		t, err := s.processOutboundOnce(ctx, input)
		if err == nil {
			return t, nil
		}
		if !errors.Is(err, models.ErrBatchVersionMismatch) {
			return nil, err
		}
		lastErr = err
	}

	return nil, fmt.Errorf("outbound transaction failed after %d optimistic locking retries: %w", s.retries, lastErr)
}

func (s *InventoryServiceImpl) processOutboundOnce(ctx context.Context, input OutboundDTO) (*models.Transaction, error) {
	t := buildTransaction(TransactionTypeOutbound, input.Reference, input.Date, input.Lines)
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.SaveTransaction(txCtx, t); err != nil {
			return err
		}

		for i := range t.Lines {
			line := &t.Lines[i]
			allocations, err := s.allocateOutboundLine(txCtx, line)
			if err != nil {
				return err
			}
			if err := s.repo.SaveBatchAllocations(txCtx, allocations); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *InventoryServiceImpl) allocateOutboundLine(ctx context.Context, line *models.TransactionLine) ([]*models.BatchAllocation, error) {
	batches, err := s.repo.AvailableBatchesForProduct(ctx, line.ProductID)
	if err != nil {
		return nil, err
	}

	remainingRequired := line.BaseQuantity
	allocations := make([]*models.BatchAllocation, 0)
	for _, batch := range batches {
		if !remainingRequired.IsPositive() {
			break
		}

		quantityToDeduct := decimal.Min(remainingRequired, batch.RemainingQuantity)
		if !quantityToDeduct.IsPositive() {
			continue
		}

		if err := s.repo.DeductBatchQuantity(ctx, batch.ID, quantityToDeduct, batch.Version); err != nil {
			return nil, err
		}

		allocations = append(allocations, &models.BatchAllocation{
			TransactionLineID: line.ID,
			InventoryBatchID:  batch.ID,
			AllocatedQuantity: quantityToDeduct,
		})
		remainingRequired = remainingRequired.Sub(quantityToDeduct)
	}

	if remainingRequired.IsPositive() {
		return nil, fmt.Errorf("%w for product %d: missing %s base units", ErrInsufficientStock, line.ProductID, remainingRequired.String())
	}

	return allocations, nil
}
