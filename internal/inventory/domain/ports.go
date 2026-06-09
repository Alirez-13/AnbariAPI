package domain

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrBatchVersionMismatch      = errors.New("inventory batch version mismatch")
	ErrInsufficientBatchQuantity = errors.New("insufficient inventory batch quantity")
)

type InventoryRepository interface {
	SaveTransaction(ctx context.Context, t *Transaction) error
	SaveInventoryBatches(ctx context.Context, batches []*InventoryBatch) error
	AvailableBatchesForProduct(ctx context.Context, productID uint) ([]*InventoryBatch, error)
	DeductBatchQuantity(ctx context.Context, batchID uint, quantityToDeduct decimal.Decimal, currentVersion uint) error
	SaveBatchAllocations(ctx context.Context, allocations []*BatchAllocation) error
}

type TransactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}
