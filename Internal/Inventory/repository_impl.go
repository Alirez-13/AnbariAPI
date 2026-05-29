package Inventory

import (
	"AnbariAPI/model"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetProduct(ctx context.Context, db *gorm.DB, id uint) (*model.Product, error) {
	var p model.Product
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, fmt.Errorf("%w: product %d", ErrProductNotFound, id)
	}
	return &p, nil
}

func (r *repository) GetBatch(ctx context.Context, db *gorm.DB, id uint, forUpdate bool) (*model.InventoryBatch, error) {
	var b model.InventoryBatch
	q := db.WithContext(ctx)
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.First(&b, id).Error; err != nil {
		return nil, fmt.Errorf("%w: batch %d", ErrBatchNotFound, id)
	}
	return &b, nil
}

func (r *repository) GetProductUnit(ctx context.Context, db *gorm.DB, productID uint, unitName string) (*model.ProductUnit, error) {
	var pu model.ProductUnit
	if err := db.WithContext(ctx).
		Where("product_id = ? AND unit_name = ?", productID, unitName).
		First(&pu).Error; err != nil {
		return nil, fmt.Errorf("%w: %q for product %d", ErrInvalidUnit, unitName, productID)
	}
	return &pu, nil
}

func (r *repository) GetAvailableBatches(ctx context.Context, productID uint) ([]model.InventoryBatch, error) {
	var batches []model.InventoryBatch
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND remaining_base_quantity > 0", productID).
		Order("entry_date ASC, id ASC").
		Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *repository) GetTransactionWithDetails(ctx context.Context, transactionID uint) (*model.Transaction, error) {
	var txn model.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Details.Product").
		Preload("Details.InventoryBatch").
		First(&txn, transactionID).Error; err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *repository) CreateTransaction(ctx context.Context, db *gorm.DB, txn *model.Transaction) error {
	return db.WithContext(ctx).Create(txn).Error
}

func (r *repository) CreateTransactionDetail(ctx context.Context, db *gorm.DB, detail *model.TransactionDetail) error {
	return db.WithContext(ctx).Create(detail).Error
}

func (r *repository) CreateInventoryBatch(ctx context.Context, db *gorm.DB, batch *model.InventoryBatch) error {
	return db.WithContext(ctx).Create(batch).Error
}

func (r *repository) UpdateProductStock(ctx context.Context, db *gorm.DB, productID uint, delta decimal.Decimal) error {
	return db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", productID).
		UpdateColumn("current_stock", gorm.Expr("current_stock + ?", delta)).
		Error
}

func (r *repository) DeductBatchStock(ctx context.Context, db *gorm.DB, batchID uint, amount decimal.Decimal) (int64, error) {
	result := db.WithContext(ctx).
		Model(&model.InventoryBatch{}).
		Where("id = ? AND remaining_base_quantity >= ?", batchID, amount).
		Update("remaining_base_quantity", gorm.Expr("remaining_base_quantity - ?", amount))
	return result.RowsAffected, result.Error
}
