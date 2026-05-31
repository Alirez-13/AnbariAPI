package service

import (
	"AnbariAPI/Internal/Inventory"
	dto2 "AnbariAPI/Internal/Inventory/dto"
	"AnbariAPI/Internal/Inventory/mapper"
	models2 "AnbariAPI/Internal/Inventory/models"
	"AnbariAPI/Internal/Inventory/repository"
	resolver2 "AnbariAPI/Internal/Inventory/resolver"
	models3 "AnbariAPI/catalog/models"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type inventoryServiceImpl struct {
	repo     repository.Repository
	resolver resolver2.ExitResolver
}

// NewInventoryService creates an InventoryService with the required dependencies.
func NewInventoryService(repo repository.Repository, resolver resolver2.ExitResolver) InventoryService {
	return &inventoryServiceImpl{
		repo:     repo,
		resolver: resolver,
	}
}

func (s *inventoryServiceImpl) GetAvailableBatches(ctx context.Context, productID uint) ([]dto2.BatchAvailabilityDTO, error) {
	batches, err := s.repo.GetAvailableBatches(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]dto2.BatchAvailabilityDTO, 0, len(batches))
	for _, b := range batches {
		result = append(result, mapper.ToBatchAvailabilityDTO(b))
	}
	return result, nil
}

func (s *inventoryServiceImpl) ProcessEntry(ctx context.Context, req dto2.EntryRequest) (*dto2.TransactionDTO, error) {
	if len(req.Lines) == 0 {
		return nil, Inventory.ErrEmptyLines
	}

	var transactionID uint

	err := s.repo.DoInTransaction(ctx, func(txRepo repository.Repository) error {
		transaction := models2.Transaction{
			TransactionType: TransactionTypeEntry,
			Reference:       req.Reference,
			Date:            req.Date,
		}
		if err := txRepo.CreateTransaction(ctx, &transaction); err != nil {
			return fmt.Errorf("failed to create entry transaction: %w", err)
		}

		transactionID = transaction.ID
		productDeltas := make(map[uint]decimal.Decimal)

		// Caches to prevent N+1 reads
		productCache := make(map[uint]*models3.Product)
		unitCache := make(map[string]decimal.Decimal)

		for _, line := range req.Lines {
			product, ok := productCache[line.ProductID]
			if !ok {
				p, err := txRepo.GetProduct(ctx, line.ProductID)
				if err != nil {
					return fmt.Errorf("process entry line for product %d: %w", line.ProductID, err)
				}
				product = p
				productCache[line.ProductID] = product
			}

			cacheKey := fmt.Sprintf("%d:%s", product.ID, line.UnitName)
			multiplier, cached := unitCache[cacheKey]
			if !cached {
				m, err := resolver2.resolveUnitMultiplier(ctx, txRepo, line.ProductID, line.UnitName, product.BaseUnit)
				if err != nil {
					return fmt.Errorf("resolve unit multiplier for product %d: %w", line.ProductID, err)
				}
				multiplier = m
				unitCache[cacheKey] = multiplier
			}

			baseQuantity := line.Quantity.Mul(multiplier)
			baseUnitPrice := line.InputUnitPrice.Div(multiplier)
			totalPrice := baseQuantity.Mul(baseUnitPrice)

			detail := models2.TransactionDetail{
				TransactionID:  transaction.ID,
				ProductID:      line.ProductID,
				UnitName:       line.UnitName,
				Multiplier:     multiplier,
				InputQuantity:  line.Quantity,
				BaseQuantity:   baseQuantity,
				InputUnitPrice: line.InputUnitPrice,
				BaseUnitPrice:  baseUnitPrice,
				TotalPrice:     totalPrice,
			}

			if err := txRepo.CreateTransactionDetail(ctx, &detail); err != nil {
				return fmt.Errorf("failed to create detail for product %d: %w", line.ProductID, err)
			}

			batch := models2.InventoryBatch{
				ProductID:             line.ProductID,
				EntryDetailID:         detail.ID,
				EntryUnitName:         line.UnitName,
				EntryUnitMultiplier:   multiplier,
				OriginalPackPrice:     line.InputUnitPrice,
				OriginalBaseUnitPrice: baseUnitPrice,
				InitialBaseQuantity:   baseQuantity,
				RemainingBaseQuantity: baseQuantity,
				EntryDate:             req.Date,
			}

			if err := txRepo.CreateInventoryBatch(ctx, &batch); err != nil {
				return fmt.Errorf("failed to create batch for product %d: %w", line.ProductID, err)
			}

			productDeltas[line.ProductID] = productDeltas[line.ProductID].Add(baseQuantity)
		}

		for productID, delta := range productDeltas {
			if err := txRepo.UpdateProductStock(ctx, productID, delta); err != nil {
				return fmt.Errorf("failed to update stock for product %d: %w", productID, err)
			}
		}

		return nil
	})

	if err != nil {
		// Wrap with the domain-specific entry failure
		return nil, fmt.Errorf("%w: %v", Inventory.ErrEntryFailed, err)
	}

	return s.loadTransactionDTO(ctx, transactionID)
}

func (s *inventoryServiceImpl) PreviewExit(ctx context.Context, req dto2.ExitRequest) (*dto2.ExitPreviewResponse, error) {
	if len(req.Lines) == 0 {
		return nil, Inventory.ErrEmptyLines
	}

	// Preview operates entirely on the read-only, non-transactional repo
	resolved, err := s.resolver.Resolve(ctx, s.repo, req.Lines, false)
	if err != nil {
		return nil, err
	}

	previewLines := make([]dto2.ExitPreviewLineDTO, 0, len(resolved))
	totalCost := decimal.Zero
	allSufficient := true

	for i, r := range resolved {
		if r.RemainingAfter.LessThan(decimal.Zero) {
			allSufficient = false
		}
		totalCost = totalCost.Add(r.LineTotal)
		previewLines = append(previewLines, mapper.ToExitPreviewLineDTO(r, req.Lines[i].UnitName))
	}

	return &dto2.ExitPreviewResponse{
		Lines:         previewLines,
		TotalCost:     totalCost,
		AllSufficient: allSufficient,
	}, nil
}

func (s *inventoryServiceImpl) ConfirmExit(ctx context.Context, req dto2.ExitRequest) (*dto2.TransactionDTO, error) {
	if len(req.Lines) == 0 {
		return nil, Inventory.ErrEmptyLines
	}

	var transactionID uint

	err := s.repo.DoInTransaction(ctx, func(txRepo repository.Repository) error {
		resolved, err := s.resolver.Resolve(ctx, txRepo, req.Lines, true)
		if err != nil {
			return err
		}

		transaction := models2.Transaction{
			TransactionType: TransactionTypeExit,
			Reference:       req.Reference,
			Date:            req.Date,
		}

		if err := txRepo.CreateTransaction(ctx, &transaction); err != nil {
			return fmt.Errorf("failed to create exit transaction: %w", err)
		}

		transactionID = transaction.ID
		productDeltas := make(map[uint]decimal.Decimal)

		for i, r := range resolved {
			batchID := r.Batch.ID
			detail := models2.TransactionDetail{
				TransactionID:    transaction.ID,
				ProductID:        r.Product.ID,
				InventoryBatchID: &batchID,
				UnitName:         req.Lines[i].UnitName,
				Multiplier:       r.Multiplier,
				InputQuantity:    r.InputQuantity,
				BaseQuantity:     r.BaseQuantity,
				InputUnitPrice:   r.InputUnitPrice,
				BaseUnitPrice:    r.BaseUnitPrice,
				TotalPrice:       r.LineTotal,
			}

			if err := txRepo.CreateTransactionDetail(ctx, &detail); err != nil {
				return fmt.Errorf("failed to create detail for batch %d: %w", batchID, err)
			}

			rowsAffected, err := txRepo.DeductBatchStock(ctx, batchID, r.BaseQuantity)
			if err != nil {
				return fmt.Errorf("failed to deduct stock for batch %d: %w", batchID, err)
			}
			if rowsAffected == 0 {
				return fmt.Errorf("%w: batch %d", Inventory.ErrConcurrentUpdate, batchID)
			}

			productDeltas[r.Product.ID] = productDeltas[r.Product.ID].Sub(r.BaseQuantity)
		}

		for productID, delta := range productDeltas {
			if err := txRepo.UpdateProductStock(ctx, productID, delta); err != nil {
				return fmt.Errorf("failed to update stock for product %d: %w", productID, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", Inventory.ErrExitFailed, err)
	}

	return s.loadTransactionDTO(ctx, transactionID)
}

func (s *inventoryServiceImpl) loadTransactionDTO(ctx context.Context, transactionID uint) (*dto2.TransactionDTO, error) {
	txn, err := s.repo.GetTransactionWithDetails(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	return mapper.ToTransactionDTO(txn), nil
}
