package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction represents main transaction information (receipts/dispatches)
type Transaction struct {
	ID          uint            `gorm:"primaryKey"`
	Type        string          `gorm:"size:20;not null;index"` // "IN" or "OUT"
	Reference   string          `gorm:"size:100;index"`         // Invoice or Document Number
	Date        time.Time       `gorm:"not null;index"`
	TotalAmount decimal.Decimal `gorm:"type:numeric(15,4);not null;default:0"` // Instant calculation

	Lines     []TransactionLine `gorm:"foreignKey:TransactionID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
