package Inventory

import (
	"AnbariAPI/model"
	"context"

	"github.com/shopspring/decimal"
)

// Repository defines the data access contract for inventory operations.
type Repository interface {
	GetProduct(ctx context.Context, id uint) (*model.Product, error)
	GetBatch(ctx context.Context, id uint) (*model.InventoryBatch, error)
	GetBatchForUpdate(ctx context.Context, id uint) (*model.InventoryBatch, error)
	GetProductUnit(ctx context.Context, productID uint, unitName string) (*model.ProductUnit, error)
	GetAvailableBatches(ctx context.Context, productID uint) ([]model.InventoryBatch, error)
	GetTransactionWithDetails(ctx context.Context, transactionID uint) (*model.Transaction, error)

	CreateTransaction(ctx context.Context, txn *model.Transaction) error
	CreateTransactionDetail(ctx context.Context, detail *model.TransactionDetail) error
	CreateInventoryBatch(ctx context.Context, batch *model.InventoryBatch) error
	UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error
	DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error)

	// DoInTransaction executes the given function within a database transaction.
	DoInTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}
