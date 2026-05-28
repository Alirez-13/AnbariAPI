package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// Entry (IN transaction) request
// ──────────────────────────────────────────────
type EntryLineRequest struct {
	ProductID      uint            `json:"product_id" binding:"required"`
	UnitName       string          `json:"unit_name" binding:"required"`             // e.g. "bucket"
	Quantity       decimal.Decimal `json:"quantity" binding:"required,gt=0"`         // e.g. 10
	InputUnitPrice decimal.Decimal `json:"input_unit_price" binding:"required,gt=0"` // e.g. $50 per bucket
}

type EntryRequest struct {
	Date      time.Time          `json:"date" binding:"required"`
	Reference string             `json:"reference,omitempty"`
	Lines     []EntryLineRequest `json:"lines" binding:"required,min=1,dive"`
}
