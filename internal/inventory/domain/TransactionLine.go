package domain

import (
	catalogdomain "AnbariAPI/internal/catalog/domain"

	"github.com/shopspring/decimal"
)

// TransactionLine represents the logical product requirement.
type TransactionLine struct {
	ID            uint                  `gorm:"primaryKey"`
	TransactionID uint                  `gorm:"not null;index"`
	ProductID     uint                  `gorm:"not null;index"`
	Product       catalogdomain.Product `gorm:"foreignKey:ProductID"`

	// Packaging Details (e.g., 2 Buckets of 25kg)
	UnitName       string          `gorm:"size:50;not null"`            // "Bucket", "Meter" , "Liter"
	UnitQuantity   decimal.Decimal `gorm:"type:numeric(15,4);not null"` // 2
	UnitMultiplier decimal.Decimal `gorm:"type:numeric(15,4);not null"` // 25
	BaseQuantity   decimal.Decimal `gorm:"type:numeric(15,4);not null"` // 50 (UnitQuantity * UnitMultiplier)

	// Financials (Entry price for IN, Exit price for OUT)
	UnitPrice decimal.Decimal `gorm:"type:numeric(15,4);not null"` // Price per UnitName
	LineTotal decimal.Decimal `gorm:"type:numeric(15,4);not null"` // UnitQuantity * UnitPrice

	// Allocations map an OUT line to specific batches.
	// (Empty for IN lines, as IN lines generate the InventoryBatch itself).
	Allocations []BatchAllocation `gorm:"foreignKey:TransactionLineID"`
}
