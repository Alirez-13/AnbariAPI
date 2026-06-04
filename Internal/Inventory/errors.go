// Package Inventory provides domain logic, domain, and repositories for managing stock.
package Inventory

import "errors"

// Domain sentinel errors for the Inventory bounded context.
var (
	ErrProductNotFound   = errors.New("product not found")
	ErrBatchNotFound     = errors.New("batch not found")
	ErrInsufficientStock = errors.New("insufficient stock in batch")
	ErrInvalidUnit       = errors.New("invalid unit for product")
	ErrEntryFailed       = errors.New("failed to process entry")
	ErrExitFailed        = errors.New("failed to process exit")
	ErrEmptyLines        = errors.New("no lines provided in request")
	ErrInvalidQuantity   = errors.New("invalid quantity: must be greater than zero")
	ErrConcurrentUpdate  = errors.New("concurrent modification detected")
)
