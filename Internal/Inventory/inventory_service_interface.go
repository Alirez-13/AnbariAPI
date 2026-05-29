package Inventory

import (
	"AnbariAPI/dto"
	"context"
)

// TransactionType constants for inventory operations.
const (
	TransactionTypeEntry = "IN"
	TransactionTypeExit  = "OUT"
)

// InventoryService provides inventory management operations for stock entries and exits.
type InventoryService interface {
	GetAvailableBatches(ctx context.Context, productID uint) ([]dto.BatchAvailabilityDTO, error)
	ProcessEntry(ctx context.Context, req dto.EntryRequest) (*dto.TransactionDTO, error)
	PreviewExit(ctx context.Context, req dto.ExitRequest) (*dto.ExitPreviewResponse, error)
	ConfirmExit(ctx context.Context, req dto.ExitRequest) (*dto.TransactionDTO, error)
}
