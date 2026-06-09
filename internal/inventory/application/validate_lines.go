package application

// # SINGLE REASON: Validate inventory transaction line input.

import "fmt"

func validateLines(lines []InventoryLineDTO) error {
	if len(lines) == 0 {
		return fmt.Errorf("%w: at least one line is required", ErrInvalidTransactionInput)
	}

	for i, line := range lines {
		if line.ProductID == 0 {
			return fmt.Errorf("%w: line %d product_id is required", ErrInvalidTransactionInput, i+1)
		}
		if line.UnitName == "" {
			return fmt.Errorf("%w: line %d unit_name is required", ErrInvalidTransactionInput, i+1)
		}
		if !line.UnitQuantity.IsPositive() {
			return fmt.Errorf("%w: line %d unit_quantity must be positive", ErrInvalidTransactionInput, i+1)
		}
		if !line.UnitMultiplier.IsPositive() {
			return fmt.Errorf("%w: line %d unit_multiplier must be positive", ErrInvalidTransactionInput, i+1)
		}
		if line.UnitPrice.IsNegative() {
			return fmt.Errorf("%w: line %d unit_price cannot be negative", ErrInvalidTransactionInput, i+1)
		}
	}

	return nil
}
