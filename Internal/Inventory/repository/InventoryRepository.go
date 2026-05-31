package repository

import (
	models2 "AnbariAPI/Internal/Inventory/models"
	models3 "AnbariAPI/catalog/models"
	"context"

	"github.com/shopspring/decimal"
)

// Repository defines the data access contract for inventory operations.
type Repository interface {
	GetProduct(ctx context.Context, id uint) (*models3.Product, error)
	GetBatch(ctx context.Context, id uint) (*models2.InventoryBatch, error)
	GetBatchForUpdate(ctx context.Context, id uint) (*models2.InventoryBatch, error)
	GetProductUnit(ctx context.Context, productID uint, unitName string) (*models3.ProductUnit, error)
	GetAvailableBatches(ctx context.Context, productID uint) ([]models2.InventoryBatch, error)
	GetTransactionWithDetails(ctx context.Context, transactionID uint) (*models2.Transaction, error)

	CreateTransaction(ctx context.Context, txn *models2.Transaction) error
	CreateTransactionDetail(ctx context.Context, detail *models2.TransactionDetail) error
	CreateInventoryBatch(ctx context.Context, batch *models2.InventoryBatch) error
	UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error
	DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error)

	// DoInTransaction executes the given function within a database transaction.
	DoInTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}
