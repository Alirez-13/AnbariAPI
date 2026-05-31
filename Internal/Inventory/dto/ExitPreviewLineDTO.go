package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// Exit preview response
// ──────────────────────────────────────────────
type ExitPreviewLineDTO struct {
	BatchID               uint            `json:"batch_id"`
	ProductID             uint            `json:"product_id"`
	ProductName           string          `json:"product_name"`
	EntryDate             time.Time       `json:"entry_date"`
	RequestedUnit         string          `json:"requested_unit"`
	RequestedQuantity     decimal.Decimal `json:"requested_quantity"`
	Multiplier            decimal.Decimal `json:"multiplier"`
	BaseQuantity          decimal.Decimal `json:"base_quantity"`
	OriginalBaseUnitPrice decimal.Decimal `json:"original_base_unit_price"`
	OriginalPackPrice     decimal.Decimal `json:"original_pack_price"` // for display: "this batch was $X per entry-unit"
	LineTotal             decimal.Decimal `json:"line_total"`          // base_quantity × base_unit_price
	RemainingBeforeExit   decimal.Decimal `json:"remaining_before_exit"`
	RemainingAfterExit    decimal.Decimal `json:"remaining_after_exit"`
	Sufficient            bool            `json:"sufficient"`
}

type ExitPreviewResponse struct {
	Lines         []ExitPreviewLineDTO `json:"lines"`
	TotalCost     decimal.Decimal      `json:"total_cost"`
	AllSufficient bool                 `json:"all_sufficient"`
}
