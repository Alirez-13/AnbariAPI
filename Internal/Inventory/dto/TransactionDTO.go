package dto

import "time"

type TransactionDTO struct {
	ID              uint                   `json:"id"`
	TransactionType string                 `json:"transaction_type"`
	Reference       string                 `json:"reference,omitempty"`
	Date            time.Time              `json:"date"`
	Details         []TransactionDetailDTO `json:"details"`
	CreatedAt       time.Time              `json:"created_at"`
}
