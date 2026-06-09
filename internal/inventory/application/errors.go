package application

// # SINGLE REASON: Define inventory application errors.

import "errors"

var (
	ErrInvalidTransactionInput = errors.New("invalid inventory transaction input")
	ErrInsufficientStock       = errors.New("insufficient stock")
)
