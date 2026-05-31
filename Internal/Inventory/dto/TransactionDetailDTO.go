package dto

import "github.com/shopspring/decimal"

// ──────────────────────────────────────────────
// Transaction response (shared for IN and OUT confirmations)
// ──────────────────────────────────────────────
type TransactionDetailDTO struct {
	ID               uint            `json:"id"`
	ProductID        uint            `json:"product_id"`
	ProductName      string          `json:"product_name"`
	InventoryBatchID *uint           `json:"inventory_batch_id,omitempty"`
	UnitName         string          `json:"unit_name"`
	Multiplier       decimal.Decimal `json:"multiplier"`
	InputQuantity    decimal.Decimal `json:"input_quantity"`
	BaseQuantity     decimal.Decimal `json:"base_quantity"`
	InputUnitPrice   decimal.Decimal `json:"input_unit_price"`
	BaseUnitPrice    decimal.Decimal `json:"base_unit_price"`
	TotalPrice       decimal.Decimal `json:"total_price"`
}
