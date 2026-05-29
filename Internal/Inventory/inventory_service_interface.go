package Inventory

import (
	"AnbariAPI/dto"
	"context"
)

type InventoryService interface {
	// ── Batch queries (for the UI popup) ──
	GetAvailableBatches(ctx context.Context, productID uint) ([]dto.BatchAvailabilityDTO, error)

	// ── Entry (IN) ──
	ProcessEntry(ctx context.Context, req dto.EntryRequest) (*dto.TransactionDTO, error)

	// ── Exit (OUT) — two-step flow ──
	PreviewExit(ctx context.Context, req dto.ExitRequest) (*dto.ExitPreviewResponse, error)
	ConfirmExit(ctx context.Context, req dto.ExitRequest) (*dto.TransactionDTO, error)
}
