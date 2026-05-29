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
// It supports a two-step exit flow where PreviewExit allows validation before ConfirmExit.
type InventoryService interface {
	// GetAvailableBatches retrieves all batches with available stock for a given product.
	// Used for UI popups to show available inventory.
	GetAvailableBatches(ctx context.Context, productID uint) ([]dto.BatchAvailabilityDTO, error)

	// ProcessEntry handles stock entry (IN) transactions.
	// Creates a transaction, details, and inventory batches atomically.
	ProcessEntry(ctx context.Context, req dto.EntryRequest) (*dto.TransactionDTO, error)

	// PreviewExit validates an exit (OUT) request and returns expected deductions.
	// Does not modify any data; useful for UI previews.
	PreviewExit(ctx context.Context, req dto.ExitRequest) (*dto.ExitPreviewResponse, error)

	// ConfirmExit commits the exit transaction after validation.
	// Consumes stock from batches and updates product totals atomically.
	ConfirmExit(ctx context.Context, req dto.ExitRequest) (*dto.TransactionDTO, error)
}
