// # SINGLE REASON: Persist inventory aggregates with GORM.
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	models "AnbariAPI/internal/inventory/domain"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GormInventoryRepository struct {
	db *gorm.DB
}

func NewGormInventoryRepository(db *gorm.DB) *GormInventoryRepository {
	return &GormInventoryRepository{db: db}
}

func (r *GormInventoryRepository) SaveTransaction(ctx context.Context, t *models.Transaction) error {
	return DBFromContext(ctx, r.db).Create(t).Error
}

func (r *GormInventoryRepository) SaveInventoryBatches(ctx context.Context, batches []*models.InventoryBatch) error {
	if len(batches) == 0 {
		return nil
	}
	return DBFromContext(ctx, r.db).Create(&batches).Error
}

func (r *GormInventoryRepository) AvailableBatchesForProduct(ctx context.Context, productID uint) ([]*models.InventoryBatch, error) {
	var batches []*models.InventoryBatch
	err := DBFromContext(ctx, r.db).
		Where("product_id = ? AND remaining_quantity > ?", productID, decimal.Zero).
		Order("entry_date ASC, id ASC").
		Find(&batches).Error
	return batches, err
}

func (r *GormInventoryRepository) DeductBatchQuantity(ctx context.Context, batchID uint, quantityToDeduct decimal.Decimal, currentVersion uint) error {
	return updateBatchQuantity(DBFromContext(ctx, r.db), batchID, quantityToDeduct, currentVersion)
}

func (r *GormInventoryRepository) SaveBatchAllocations(ctx context.Context, allocations []*models.BatchAllocation) error {
	if len(allocations) == 0 {
		return nil
	}
	return DBFromContext(ctx, r.db).Create(&allocations).Error
}

func (r *GormInventoryRepository) CreateTransaction(tx *gorm.DB, t *models.Transaction) error {
	return tx.Create(t).Error
}

func (r *GormInventoryRepository) CreateInventoryBatches(tx *gorm.DB, batches []*models.InventoryBatch) error {
	if len(batches) == 0 {
		return nil
	}
	return tx.Create(&batches).Error
}

func (r *GormInventoryRepository) GetAvailableBatchesForProduct(ctx context.Context, productID uint) ([]*models.InventoryBatch, error) {
	return r.AvailableBatchesForProduct(ctx, productID)
}

func (r *GormInventoryRepository) UpdateBatchQuantity(tx *gorm.DB, batchID uint, quantityToDeduct decimal.Decimal, currentVersion uint) error {
	return updateBatchQuantity(tx, batchID, quantityToDeduct, currentVersion)
}

func (r *GormInventoryRepository) CreateBatchAllocations(tx *gorm.DB, allocations []*models.BatchAllocation) error {
	if len(allocations) == 0 {
		return nil
	}
	return tx.Create(&allocations).Error
}

func updateBatchQuantity(db *gorm.DB, batchID uint, quantityToDeduct decimal.Decimal, currentVersion uint) error {
	if !quantityToDeduct.IsPositive() {
		return fmt.Errorf("quantity to deduct must be positive")
	}

	result := db.Model(&models.InventoryBatch{}).
		Where("id = ? AND version = ? AND remaining_quantity >= ?", batchID, currentVersion, quantityToDeduct).
		Updates(map[string]any{
			"remaining_quantity": gorm.Expr("remaining_quantity - ?", quantityToDeduct),
			"version":            gorm.Expr("version + 1"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var batch models.InventoryBatch
	err := db.Select("id", "version", "remaining_quantity").First(&batch, batchID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrInsufficientBatchQuantity
	}
	if err != nil {
		return err
	}
	if batch.Version != currentVersion {
		return models.ErrBatchVersionMismatch
	}
	return models.ErrInsufficientBatchQuantity
}
