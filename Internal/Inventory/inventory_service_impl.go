package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type inventoryService struct {
	repo     Repository
	resolver ExitResolver
}

// NewInventoryService creates an InventoryService with the required dependencies.
// TODO: Update injection graph. Replace `NewInventoryService(db)` with `NewInventoryService(repo, resolver)`.
func NewInventoryService(repo Repository, resolver ExitResolver) InventoryService {
	return &inventoryService{
		repo:     repo,
		resolver: resolver,
	}
}

func (s *inventoryService) GetAvailableBatches(ctx context.Context, productID uint) ([]dto.BatchAvailabilityDTO, error) {
	batches, err := s.repo.GetAvailableBatches(ctx, productID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.BatchAvailabilityDTO, 0, len(batches))
	for _, b := range batches {
		result = append(result, ToBatchAvailabilityDTO(b))
	}
	return result, nil
}

func (s *inventoryService) ProcessEntry(ctx context.Context, req dto.EntryRequest) (*dto.TransactionDTO, error) {
	if len(req.Lines) == 0 {
		return nil, ErrEmptyLines
	}

	var transactionID uint

	err := s.repo.DoInTransaction(ctx, func(txRepo Repository) error {
		transaction := model.Transaction{
			TransactionType: TransactionTypeEntry,
			Reference:       req.Reference,
			Date:            req.Date,
		}
		if err := txRepo.CreateTransaction(ctx, &transaction); err != nil {
			return fmt.Errorf("%w: create transaction: %v", ErrEntryFailed, err)
		}

		transactionID = transaction.ID
		productDeltas := make(map[uint]decimal.Decimal)

		for _, line := range req.Lines {
			product, err := txRepo.GetProduct(ctx, line.ProductID)
			if err != nil {
				return fmt.Errorf("process entry line for product %d: %w", line.ProductID, err)
			}

			multiplier, err := resolveUnitMultiplier(ctx, txRepo, line.ProductID, line.UnitName, product.BaseUnit)
			if err != nil {
				return fmt.Errorf("resolve unit multiplier for product %d: %w", line.ProductID, err)
			}

			baseQuantity := line.Quantity.Mul(multiplier)
			baseUnitPrice := line.InputUnitPrice.Div(multiplier)
			totalPrice := baseQuantity.Mul(baseUnitPrice)

			detail := model.TransactionDetail{
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
				return fmt.Errorf("%w: create detail: %v", ErrEntryFailed, err)
			}

			batch := model.InventoryBatch{
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
				return fmt.Errorf("%w: create batch: %v", ErrEntryFailed, err)
			}

			productDeltas[line.ProductID] = productDeltas[line.ProductID].Add(baseQuantity)
		}

		for productID, delta := range productDeltas {
			if err := txRepo.UpdateProductStock(ctx, productID, delta); err != nil {
				return fmt.Errorf("%w: update product stock: %v", ErrEntryFailed, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.loadTransactionDTO(ctx, transactionID)
}

func (s *inventoryService) PreviewExit(ctx context.Context, req dto.ExitRequest) (*dto.ExitPreviewResponse, error) {
	if len(req.Lines) == 0 {
		return nil, ErrEmptyLines
	}

	// Preview operates entirely on the read-only, non-transactional repo
	resolved, err := s.resolver.Resolve(ctx, s.repo, req.Lines, false)
	if err != nil {
		return nil, err
	}

	previewLines := make([]dto.ExitPreviewLineDTO, 0, len(resolved))
	totalCost := decimal.Zero
	allSufficient := true

	for i, r := range resolved {
		if r.RemainingAfter.LessThan(decimal.Zero) {
			allSufficient = false
		}
		totalCost = totalCost.Add(r.LineTotal)
		previewLines = append(previewLines, ToExitPreviewLineDTO(r, req.Lines[i].UnitName))
	}

	return &dto.ExitPreviewResponse{
		Lines:         previewLines,
		TotalCost:     totalCost,
		AllSufficient: allSufficient,
	}, nil
}

func (s *inventoryService) ConfirmExit(ctx context.Context, req dto.ExitRequest) (*dto.TransactionDTO, error) {
	if len(req.Lines) == 0 {
		return nil, ErrEmptyLines
	}

	var transactionID uint

	err := s.repo.DoInTransaction(ctx, func(txRepo Repository) error {
		// Pass the transactional repository to the resolver
		resolved, err := s.resolver.Resolve(ctx, txRepo, req.Lines, true)
		if err != nil {
			return err
		}

		transaction := model.Transaction{
			TransactionType: TransactionTypeExit,
			Reference:       req.Reference,
			Date:            req.Date,
		}
		if err := txRepo.CreateTransaction(ctx, &transaction); err != nil {
			return fmt.Errorf("%w: create transaction: %v", ErrExitFailed, err)
		}

		transactionID = transaction.ID
		productDeltas := make(map[uint]decimal.Decimal)

		for i, r := range resolved {
			batchID := r.Batch.ID
			detail := model.TransactionDetail{
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
				return fmt.Errorf("%w: create detail: %v", ErrExitFailed, err)
			}

			rowsAffected, err := txRepo.DeductBatchStock(ctx, batchID, r.BaseQuantity)
			if err != nil {
				return fmt.Errorf("%w: update batch %d: %v", ErrExitFailed, batchID, err)
			}
			if rowsAffected == 0 {
				return fmt.Errorf("%w: batch %d — concurrent modification detected", ErrInsufficientStock, batchID)
			}

			productDeltas[r.Product.ID] = productDeltas[r.Product.ID].Sub(r.BaseQuantity)
		}

		for productID, delta := range productDeltas {
			if err := txRepo.UpdateProductStock(ctx, productID, delta); err != nil {
				return fmt.Errorf("%w: update product stock: %v", ErrExitFailed, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.loadTransactionDTO(ctx, transactionID)
}

func (s *inventoryService) loadTransactionDTO(ctx context.Context, transactionID uint) (*dto.TransactionDTO, error) {
	txn, err := s.repo.GetTransactionWithDetails(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	return ToTransactionDTO(txn), nil
}
