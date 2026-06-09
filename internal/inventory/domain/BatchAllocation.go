package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// BatchAllocation resolves the Many-to-Many between Lines and Batches for OUTs.
type BatchAllocation struct {
	ID                uint `gorm:"primaryKey"`
	TransactionLineID uint `gorm:"not null;index"` // The "OUT" line consuming stock
	InventoryBatchID  uint `gorm:"not null;index"` // The specific batch being consumed

	AllocatedQuantity decimal.Decimal `gorm:"type:numeric(15,4);not null"` // Base units consumed

	CreatedAt time.Time
}
