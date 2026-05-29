package Inventory

import "fmt"

var (
	ErrProductNotFound   = fmt.Errorf("product not found")
	ErrBatchNotFound     = fmt.Errorf("batch not found")
	ErrInsufficientStock = fmt.Errorf("insufficient stock in batch")
	ErrInvalidUnit       = fmt.Errorf("invalid unit for product")
	ErrEntryFailed       = fmt.Errorf("failed to process entry")
	ErrExitFailed        = fmt.Errorf("failed to process exit")
)
