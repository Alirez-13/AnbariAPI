package service

import (
	dto2 "AnbariAPI/Internal/Inventory/dto"
	"context"
)

// TransactionType constants for inventory operations.
const (
	TransactionTypeEntry = "IN"
	TransactionTypeExit  = "OUT"
)

// InventoryService provides inventory management operations for stock entries and exits.
type InventoryService interface {
	GetAvailableBatches(ctx context.Context, productID uint) ([]dto2.BatchAvailabilityDTO, error)
	ProcessEntry(ctx context.Context, req dto2.EntryRequest) (*dto2.TransactionDTO, error)
	PreviewExit(ctx context.Context, req dto2.ExitRequest) (*dto2.ExitPreviewResponse, error)
	ConfirmExit(ctx context.Context, req dto2.ExitRequest) (*dto2.TransactionDTO, error)
}
