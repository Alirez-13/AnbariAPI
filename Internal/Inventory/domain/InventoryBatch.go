package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// InventoryBatch — KEY CHANGES
//
// Added snapshotted entry fields so every batch
// is self-describing without joining back.
// ──────────────────────────────────────────────
type InventoryBatch struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProductID uint `gorm:"not null;index" json:"product_id"`

	// Audit link back to the IN TransactionDetail that created this batch.
	// Immutable after creation.
	EntryDetailID uint `gorm:"not null;index" json:"entry_detail_id"`

	// ── Snapshotted entry context (immutable) ───────────────────
	// These are frozen at entry time so that even if ProductUnit
	// definitions change later, the batch record remains accurate.
	EntryUnitName         string          `gorm:"size:50;not null" json:"entry_unit_name"`                     // e.g. "bucket"
	EntryUnitMultiplier   decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"entry_unit_multiplier"`    // e.g. 25
	OriginalPackPrice     decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"original_pack_price"`      // e.g. $50.00 per bucket
	OriginalBaseUnitPrice decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"original_base_unit_price"` // e.g. $2.00 per kg
	// ─────────────────────────────────────────────────────────────

	InitialBaseQuantity   decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"initial_base_quantity"`
	RemainingBaseQuantity decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"remaining_base_quantity"`
	EntryDate             time.Time       `gorm:"not null;index" json:"entry_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsExhausted returns true when the batch has no remaining stock.
func (b *InventoryBatch) IsExhausted() bool {
	return b.RemainingBaseQuantity.LessThanOrEqual(decimal.Zero)
}
