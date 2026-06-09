package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type InventoryBatch struct {
	ID                uint `gorm:"primaryKey"`
	TransactionLineID uint `gorm:"not null;uniqueIndex"` // The "IN" line that created this stock
	ProductID         uint `gorm:"not null;index"`

	// Tracking in Base Units (e.g., kg, liters)
	InitialQuantity   decimal.Decimal `gorm:"type:numeric(15,4);not null"`
	RemainingQuantity decimal.Decimal `gorm:"type:numeric(15,4);not null"`

	// Cost stored per Base Unit for standard valuation
	EntryUnitCost decimal.Decimal `gorm:"type:numeric(15,4);not null"`

	EntryDate time.Time `gorm:"not null;index"`
	Version   uint      `gorm:"not null;default:1"` // Optimistic locking

	CreatedAt time.Time
	UpdatedAt time.Time
}
