package Inventory

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrBatchNotFound     = errors.New("batch not found")
	ErrInsufficientStock = errors.New("insufficient stock in batch")
	ErrInvalidUnit       = errors.New("invalid unit for product")
	ErrEntryFailed       = errors.New("failed to process entry")
	ErrExitFailed        = errors.New("failed to process exit")
	ErrEmptyLines        = errors.New("no lines provided in request")
	ErrInvalidQuantity   = errors.New("invalid quantity: must be greater than zero")
)