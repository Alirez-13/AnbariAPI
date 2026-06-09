package application

// # SINGLE REASON: Define inventory use case input DTOs.

import (
	"time"

	"github.com/shopspring/decimal"
)

type InboundDTO struct {
	Reference string
	Date      time.Time
	Lines     []InventoryLineDTO
}

type OutboundDTO struct {
	Reference string
	Date      time.Time
	Lines     []InventoryLineDTO
}

type InventoryLineDTO struct {
	ProductID      uint
	UnitName       string
	UnitQuantity   decimal.Decimal
	UnitMultiplier decimal.Decimal
	UnitPrice      decimal.Decimal
}
