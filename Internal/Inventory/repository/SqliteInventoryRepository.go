package repository

import (
	"AnbariAPI/Internal/Inventory"
	domain2 "AnbariAPI/Internal/Inventory/domain"
	"AnbariAPI/Internal/catalog/domain"
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new database-backed inventory repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DoInTransaction(ctx context.Context, fn func(txRepo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repository{db: tx})
	})
}

func (r *repository) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	var p domain.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: ID %d", Inventory.ErrProductNotFound, id)
		}
		return nil, fmt.Errorf("failed to get product %d: %w", id, err)
	}
	return &p, nil
}

func (r *repository) GetBatch(ctx context.Context, id uint) (*domain2.InventoryBatch, error) {
	var b domain2.InventoryBatch
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: ID %d", Inventory.ErrBatchNotFound, id)
		}
		return nil, fmt.Errorf("failed to get batch %d: %w", id, err)
	}
	return &b, nil
}

func (r *repository) GetBatchForUpdate(ctx context.Context, id uint) (*domain2.InventoryBatch, error) {
	var b domain2.InventoryBatch
	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: ID %d", Inventory.ErrBatchNotFound, id)
		}
		return nil, fmt.Errorf("failed to get batch for update %d: %w", id, err)
	}
	return &b, nil
}

func (r *repository) GetAvailableBatches(ctx context.Context, productID uint) ([]domain2.InventoryBatch, error) {
	var batches []domain2.InventoryBatch
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND remaining_base_quantity > 0", productID).
		Order("entry_date ASC, id ASC").
		Find(&batches).Error; err != nil {
		return nil, fmt.Errorf("failed to get available batches for product %d: %w", productID, err)
	}
	return batches, nil
}

func (r *repository) GetTransactionWithDetails(ctx context.Context, transactionID uint) (*domain2.Transaction, error) {
	var txn domain2.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Details.Product").
		Preload("Details.InventoryBatch").
		First(&txn, transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction %d not found: %w", transactionID, err)
		}
		return nil, fmt.Errorf("failed to get transaction %d with details: %w", transactionID, err)
	}
	return &txn, nil
}

func (r *repository) CreateTransaction(ctx context.Context, txn *domain2.Transaction) error {
	if err := r.db.WithContext(ctx).Create(txn).Error; err != nil {
		return fmt.Errorf("create transaction error: %w", err)
	}
	return nil
}

func (r *repository) CreateTransactionDetail(ctx context.Context, detail *domain2.TransactionDetail) error {
	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return fmt.Errorf("create transaction detail for product %d: %w", detail.ProductID, err)
	}
	return nil
}

func (r *repository) CreateInventoryBatch(ctx context.Context, batch *domain2.InventoryBatch) error {
	if err := r.db.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("create inventory batch for product %d: %w", batch.ProductID, err)
	}
	return nil
}

func (r *repository) UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		UpdateColumn("current_stock", gorm.Expr("current_stock + ?", delta))

	if result.Error != nil {
		return fmt.Errorf("update stock for product %d: %w", productID, result.Error)
	}
	return nil
}

func (r *repository) DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&domain2.InventoryBatch{}).
		Where("id = ? AND remaining_base_quantity >= ?", batchID, amount).
		UpdateColumn("remaining_base_quantity", gorm.Expr("remaining_base_quantity - ?", amount))

	if result.Error != nil {
		return 0, fmt.Errorf("deduct stock from batch %d: %w", batchID, result.Error)
	}
	return result.RowsAffected, nil
}
