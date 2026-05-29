package Inventory

import (
	"AnbariAPI/model"
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
	// GORM's Transaction automatically handles panics and rollbacks.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Provide a new repository instance bound strictly to this transaction
		return fn(&repository{db: tx})
	})
}

func (r *repository) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	var p model.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: product %d", ErrProductNotFound, id)
		}
		return nil, fmt.Errorf("failed to get product %d: %w", id, err)
	}
	return &p, nil
}

func (r *repository) GetBatch(ctx context.Context, id uint, forUpdate bool) (*model.InventoryBatch, error) {
	var b model.InventoryBatch
	q := r.db.WithContext(ctx)
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: batch %d", ErrBatchNotFound, id)
		}
		return nil, fmt.Errorf("failed to get batch %d: %w", id, err)
	}
	return &b, nil
}

func (r *repository) GetProductUnit(ctx context.Context, productID uint, unitName string) (*model.ProductUnit, error) {
	var pu model.ProductUnit
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND unit_name = ?", productID, unitName).
		First(&pu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %q for product %d", ErrInvalidUnit, unitName, productID)
		}
		return nil, fmt.Errorf("failed to get product unit %q for product %d: %w", unitName, productID, err)
	}
	return &pu, nil
}

func (r *repository) GetAvailableBatches(ctx context.Context, productID uint) ([]model.InventoryBatch, error) {
	var batches []model.InventoryBatch
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND remaining_base_quantity > 0", productID).
		Order("entry_date ASC, id ASC").
		Find(&batches).Error; err != nil {
		return nil, fmt.Errorf("failed to get available batches for product %d: %w", productID, err)
	}
	return batches, nil
}

func (r *repository) GetTransactionWithDetails(ctx context.Context, transactionID uint) (*model.Transaction, error) {
	var txn model.Transaction
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

func (r *repository) CreateTransaction(ctx context.Context, txn *model.Transaction) error {
	if err := r.db.WithContext(ctx).Create(txn).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *repository) CreateTransactionDetail(ctx context.Context, detail *model.TransactionDetail) error {
	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return fmt.Errorf("failed to create transaction detail for product %d: %w", detail.ProductID, err)
	}
	return nil
}

func (r *repository) CreateInventoryBatch(ctx context.Context, batch *model.InventoryBatch) error {
	if err := r.db.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("failed to create inventory batch for product %d: %w", batch.ProductID, err)
	}
	return nil
}

func (r *repository) UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error {
	result := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", productID).
		UpdateColumn("current_stock", gorm.Expr("current_stock + ?", delta))
	if result.Error != nil {
		return fmt.Errorf("failed to update stock for product %d: %w", productID, result.Error)
	}
	return nil
}

func (r *repository) DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&model.InventoryBatch{}).
		Where("id = ? AND remaining_base_quantity >= ?", batchID, amount).
		Update("remaining_base_quantity", gorm.Expr("remaining_base_quantity - ?", amount))
	if result.Error != nil {
		return 0, fmt.Errorf("failed to deduct stock from batch %d: %w", batchID, result.Error)
	}
	return result.RowsAffected, nil
}
