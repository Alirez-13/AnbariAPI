package model

import (
	"time"
)

// Transaction represents main transaction information (receipts/dispatches)
type Transaction struct {
	ID              uint                `gorm:"primaryKey" json:"id"`
	TransactionType string              `gorm:"size:20;not null" json:"transaction_type"` // "IN" or "OUT"
	Reference       string              `gorm:"size:100" json:"reference,omitempty"`      // invoice/delivery note
	Date            time.Time           `gorm:"not null;index" json:"date"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
}
