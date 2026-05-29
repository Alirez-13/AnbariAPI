package Inventory

import (
	"AnbariAPI/model"
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Repository defines the data access contract for inventory operations.
// All methods accept a db parameter to allow transaction control at the service layer.
type Repository interface {
	GetProduct(ctx context.Context, db *gorm.DB, id uint) (*model.Product, error)
	GetBatch(ctx context.Context, db *gorm.DB, id uint, forUpdate bool) (*model.InventoryBatch, error)
	GetProductUnit(ctx context.Context, db *gorm.DB, productID uint, unitName string) (*model.ProductUnit, error)
	GetAvailableBatches(ctx context.Context, productID uint) ([]model.InventoryBatch, error)
	GetTransactionWithDetails(ctx context.Context, transactionID uint) (*model.Transaction, error)

	CreateTransaction(ctx context.Context, db *gorm.DB, txn *model.Transaction) error
	CreateTransactionDetail(ctx context.Context, db *gorm.DB, detail *model.TransactionDetail) error
	CreateInventoryBatch(ctx context.Context, db *gorm.DB, batch *model.InventoryBatch) error
	UpdateProductStock(ctx context.Context, db *gorm.DB, productID uint, delta decimal.Decimal) error
	DeductBatchStock(ctx context.Context, db *gorm.DB, batchID uint, amount decimal.Decimal) (int64, error)
}