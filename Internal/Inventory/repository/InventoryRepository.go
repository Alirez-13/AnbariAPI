package repository

import (
	domain2 "AnbariAPI/Internal/Inventory/domain"
	"AnbariAPI/Internal/catalog/domain"
	"context"

	"github.com/shopspring/decimal"
)

// Repository defines the data access contract for inventory operations.
type Repository interface {
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	GetBatch(ctx context.Context, id uint) (*domain2.InventoryBatch, error)
	GetBatchForUpdate(ctx context.Context, id uint) (*domain2.InventoryBatch, error)
	GetProductUnit(ctx context.Context, productID uint, unitName string) (*domain.ProductUnit, error)
	GetAvailableBatches(ctx context.Context, productID uint) ([]domain2.InventoryBatch, error)
	GetTransactionWithDetails(ctx context.Context, transactionID uint) (*domain2.Transaction, error)

	CreateTransaction(ctx context.Context, txn *domain2.Transaction) error
	CreateTransactionDetail(ctx context.Context, detail *domain2.TransactionDetail) error
	CreateInventoryBatch(ctx context.Context, batch *domain2.InventoryBatch) error
	UpdateProductStock(ctx context.Context, productID uint, delta decimal.Decimal) error
	DeductBatchStock(ctx context.Context, batchID uint, amount decimal.Decimal) (int64, error)

	// DoInTransaction executes the given function within a database transaction.
	DoInTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}
