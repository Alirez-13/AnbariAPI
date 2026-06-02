package repository

import (
	"AnbariAPI/shared/models"
	"context"

	"github.com/shopspring/decimal"
)

// Repository defines the data access contract for inventory operations.
type Repository interface {
	GetProduct(ctx context.Context, id uint) (*models.Product, error)
	GetBatch(ctx context.Context, id uint) (*models.InventoryBatch, error)
	GetBatchForUpdate(ctx context.Context, id uint) (*models.InventoryBatch, error)
	GetProductUnit(ctx context.Context, productID uint, unitName string) (*models.ProductUnit, error)
	GetAvailableBatches(ctx context.Context, productID uint) ([]models.InventoryBatch, error)
	GetTransactionWithDetails(ctx context.Context, transactionID uint) (*models.Transaction, error)

	CreateTransaction(ctx context.Context, txn *models.Transaction) error
	CreateTransactionDetail(ctx context.Context, detail *models.TransactionDetail) error
	CreateInventoryBatch(ctx context.Context, batch *models.InventoryBatch) error
	UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error
	DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error)

	// DoInTransaction executes the given function within a database transaction.
	DoInTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}
