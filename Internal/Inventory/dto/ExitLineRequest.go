package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// Exit (OUT transaction) — shared by Preview & Confirm
// ──────────────────────────────────────────────
type ExitLineRequest struct {
	BatchID  uint            `json:"batch_id" binding:"required"`
	UnitName string          `json:"unit_name" binding:"required"` // unit to exit in (can differ from entry unit)
	Quantity decimal.Decimal `json:"quantity" binding:"required,gt=0"`
}

type ExitRequest struct {
	Date      time.Time         `json:"date" binding:"required"`
	Reference string            `json:"reference,omitempty"`
	Lines     []ExitLineRequest `json:"lines" binding:"required,min=1,dive"`
}
