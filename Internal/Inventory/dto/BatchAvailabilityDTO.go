package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// Batch availability — returned by GET /products/:id/batches
// ──────────────────────────────────────────────
type BatchAvailabilityDTO struct {
	BatchID               uint            `json:"batch_id"`
	EntryDate             time.Time       `json:"entry_date"`
	EntryUnitName         string          `json:"entry_unit_name"`          // "bucket"
	EntryUnitMultiplier   decimal.Decimal `json:"entry_unit_multiplier"`    // 25
	OriginalPackPrice     decimal.Decimal `json:"original_pack_price"`      // $50.00 per bucket
	OriginalBaseUnitPrice decimal.Decimal `json:"original_base_unit_price"` // $2.00 per kg
	RemainingBaseQuantity decimal.Decimal `json:"remaining_base_quantity"`  // 150 kg left
}
