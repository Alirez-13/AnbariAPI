package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type inventoryService struct {
	db       *gorm.DB
	repo     Repository
	resolver *exitResolver
}

func NewInventoryService(db *gorm.DB) InventoryService {
	repo := NewRepository(db)
	return &inventoryService{
		db:       db,
		repo:     repo,
		resolver: newExitResolver(repo),
	}
}

func (s *inventoryService) GetAvailableBatches(ctx context.Context, productID uint) ([]dto.BatchAvailabilityDTO, error) {
	batches, err := s.repo.GetAvailableBatches(ctx, productID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.BatchAvailabilityDTO, 0, len(batches))
	for _, b := range batches {
		result = append(result, toBatchAvailabilityDTO(b))
	}
	return result, nil
}

func (s *inventoryService) ProcessEntry(ctx context.Context, req dto.EntryRequest) (*dto.TransactionDTO, error) {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := model.Transaction{
		TransactionType: "IN",
		Reference:       req.Reference,
		Date:            req.Date,
	}
	if err := s.repo.CreateTransaction(ctx, tx, &transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w: create transaction: %v", ErrEntryFailed, err)
	}

	productDeltas := make(map[uint]decimal.Decimal)

	for _, line := range req.Lines {
		product, err := s.repo.GetProduct(ctx, tx, line.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		multiplier, err := resolveUnitMultiplier(ctx, tx, s.repo, line.ProductID, line.UnitName, product.BaseUnit)
		if err != nil {
			tx.Rollback()
			return nil, err
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
		if err := s.repo.CreateTransactionDetail(ctx, tx, &detail); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: create detail: %v", ErrEntryFailed, err)
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
		if err := s.repo.CreateInventoryBatch(ctx, tx, &batch); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: create batch: %v", ErrEntryFailed, err)
		}

		productDeltas[line.ProductID] = productDeltas[line.ProductID].Add(baseQuantity)
	}

	for productID, delta := range productDeltas {
		if err := s.repo.UpdateProductStock(ctx, tx, productID, delta); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: update product stock: %v", ErrEntryFailed, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("%w: commit: %v", ErrEntryFailed, err)
	}

	return s.loadTransactionDTO(ctx, transaction.ID)
}

func (s *inventoryService) PreviewExit(ctx context.Context, req dto.ExitRequest) (*dto.ExitPreviewResponse, error) {
	resolved, err := s.resolver.resolve(ctx, s.db, req.Lines, false)
	if err != nil {
		return nil, err
	}

	previewLines := make([]dto.ExitPreviewLineDTO, 0, len(resolved))
	totalCost := decimal.Zero
	allSufficient := true

	for i, r := range resolved {
		if r.remainingAfter.LessThan(decimal.Zero) {
			allSufficient = false
		}
		totalCost = totalCost.Add(r.lineTotal)
		previewLines = append(previewLines, toExitPreviewLineDTO(r, req.Lines[i].UnitName))
	}

	return &dto.ExitPreviewResponse{
		Lines:         previewLines,
		TotalCost:     totalCost,
		AllSufficient: allSufficient,
	}, nil
}

func (s *inventoryService) ConfirmExit(ctx context.Context, req dto.ExitRequest) (*dto.TransactionDTO, error) {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	resolved, err := s.resolver.resolve(ctx, tx, req.Lines, true)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := model.Transaction{
		TransactionType: "OUT",
		Reference:       req.Reference,
		Date:            req.Date,
	}
	if err := s.repo.CreateTransaction(ctx, tx, &transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w: create transaction: %v", ErrExitFailed, err)
	}

	productDeltas := make(map[uint]decimal.Decimal)

	for i, r := range resolved {
		batchID := r.batch.ID
		detail := model.TransactionDetail{
			TransactionID:    transaction.ID,
			ProductID:        r.product.ID,
			InventoryBatchID: &batchID,
			UnitName:         req.Lines[i].UnitName,
			Multiplier:       r.multiplier,
			InputQuantity:    r.inputQuantity,
			BaseQuantity:     r.baseQuantity,
			InputUnitPrice:   r.inputUnitPrice,
			BaseUnitPrice:    r.baseUnitPrice,
			TotalPrice:       r.lineTotal,
		}
		if err := s.repo.CreateTransactionDetail(ctx, tx, &detail); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: create detail: %v", ErrExitFailed, err)
		}

		rowsAffected, err := s.repo.DeductBatchStock(ctx, tx, batchID, r.baseQuantity)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: update batch %d: %v", ErrExitFailed, batchID, err)
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("%w: batch %d — concurrent modification detected or insufficient stock", ErrInsufficientStock, batchID)
		}

		productDeltas[r.product.ID] = productDeltas[r.product.ID].Sub(r.baseQuantity)
	}

	for productID, delta := range productDeltas {
		if err := s.repo.UpdateProductStock(ctx, tx, productID, delta); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%w: update product stock: %v", ErrExitFailed, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("%w: commit: %v", ErrExitFailed, err)
	}

	return s.loadTransactionDTO(ctx, transaction.ID)
}

func (s *inventoryService) loadTransactionDTO(ctx context.Context, transactionID uint) (*dto.TransactionDTO, error) {
	txn, err := s.repo.GetTransactionWithDetails(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	return toTransactionDTO(txn), nil
}
