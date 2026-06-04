package domain

import (
	"github.com/shopspring/decimal"
)

// TransactionDetail
// ──────────────────────────────────────────────
// For IN:  InputUnitPrice = price the user typed (per entry unit)
//
//	BaseUnitPrice  = InputUnitPrice / Multiplier
//
// For OUT: InputUnitPrice = BaseUnitPrice * Multiplier (price per exit unit)
//
//	BaseUnitPrice  = batch.OriginalBaseUnitPrice
//
// ──────────────────────────────────────────────
type TransactionDetail struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	TransactionID uint `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint `gorm:"not null;index" json:"product_id"`
	//Product          domain.Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	InventoryBatchID *uint           `gorm:"index" json:"inventory_batch_id,omitempty"` // NULL for IN, required for OUT
	InventoryBatch   *InventoryBatch `gorm:"foreignKey:InventoryBatchID" json:"inventory_batch,omitempty"`

	UnitName       string          `gorm:"size:50;not null" json:"unit_name"`
	Multiplier     decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"multiplier"`
	InputQuantity  decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"input_quantity"`
	BaseQuantity   decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"base_quantity"`
	InputUnitPrice decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"input_unit_price"` // NEW
	BaseUnitPrice  decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"base_unit_price"`
	TotalPrice     decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"total_price"`
}
