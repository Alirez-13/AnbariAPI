package application

// # SINGLE REASON: Build inventory transaction aggregates from input lines.

import (
	"time"

	models "AnbariAPI/internal/inventory/domain"
)

const (
	TransactionTypeInbound  = "IN"
	TransactionTypeOutbound = "OUT"
)

func buildTransaction(transactionType, reference string, date time.Time, lines []InventoryLineDTO) *models.Transaction {
	if date.IsZero() {
		date = time.Now().UTC()
	}

	t := &models.Transaction{
		Type:      transactionType,
		Reference: reference,
		Date:      date,
		Lines:     make([]models.TransactionLine, 0, len(lines)),
	}

	for _, inputLine := range lines {
		baseQuantity := inputLine.UnitQuantity.Mul(inputLine.UnitMultiplier)
		lineTotal := inputLine.UnitQuantity.Mul(inputLine.UnitPrice)
		t.TotalAmount = t.TotalAmount.Add(lineTotal)
		t.Lines = append(t.Lines, models.TransactionLine{
			ProductID:      inputLine.ProductID,
			UnitName:       inputLine.UnitName,
			UnitQuantity:   inputLine.UnitQuantity,
			UnitMultiplier: inputLine.UnitMultiplier,
			BaseQuantity:   baseQuantity,
			UnitPrice:      inputLine.UnitPrice,
			LineTotal:      lineTotal,
		})
	}

	return t
}
