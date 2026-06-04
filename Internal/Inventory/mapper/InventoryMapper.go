package mapper

import (
	"AnbariAPI/Internal/Inventory/domain"
	dto2 "AnbariAPI/Internal/Inventory/dto"
	"AnbariAPI/Internal/Inventory/resolver"

	"github.com/shopspring/decimal"
)

// ToBatchAvailabilityDTO maps a database batch domain to an availability DTO.
func ToBatchAvailabilityDTO(b domain.InventoryBatch) dto2.BatchAvailabilityDTO {
	return dto2.BatchAvailabilityDTO{
		BatchID:               b.ID,
		EntryDate:             b.EntryDate,
		EntryUnitName:         b.EntryUnitName,
		EntryUnitMultiplier:   b.EntryUnitMultiplier,
		OriginalPackPrice:     b.OriginalPackPrice,
		OriginalBaseUnitPrice: b.OriginalBaseUnitPrice,
		RemainingBaseQuantity: b.RemainingBaseQuantity,
	}
}

// ToTransactionDTO maps a database transaction domain and its details to a DTO.
func ToTransactionDTO(txn *domain.Transaction) *dto2.TransactionDTO {
	if txn == nil {
		return nil
	}

	details := make([]dto2.TransactionDetailDTO, 0, len(txn.Details))
	for _, d := range txn.Details {
		productName := ""
		// Safe pointer checks
		if d.Product.ID != 0 {
			productName = d.Product.Name
		}

		details = append(details, dto2.TransactionDetailDTO{
			ID:               d.ID,
			ProductID:        d.ProductID,
			ProductName:      productName,
			InventoryBatchID: d.InventoryBatchID,
			UnitName:         d.UnitName,
			Multiplier:       d.Multiplier,
			InputQuantity:    d.InputQuantity,
			BaseQuantity:     d.BaseQuantity,
			InputUnitPrice:   d.InputUnitPrice,
			BaseUnitPrice:    d.BaseUnitPrice,
			TotalPrice:       d.TotalPrice,
		})
	}

	return &dto2.TransactionDTO{
		ID:              txn.ID,
		TransactionType: txn.TransactionType,
		Reference:       txn.Reference,
		Date:            txn.Date,
		Details:         details,
		CreatedAt:       txn.CreatedAt,
	}
}

// ToExitPreviewLineDTO converts a resolved domain line into a UI-friendly preview object.
func ToExitPreviewLineDTO(r resolver.ResolvedExitLine, requestedUnit string) dto2.ExitPreviewLineDTO {
	return dto2.ExitPreviewLineDTO{
		BatchID:               r.Batch.ID,
		ProductID:             r.Product.ID,
		ProductName:           r.Product.Name,
		EntryDate:             r.Batch.EntryDate,
		RequestedUnit:         requestedUnit,
		RequestedQuantity:     r.InputQuantity,
		Multiplier:            r.Multiplier,
		BaseQuantity:          r.BaseQuantity,
		OriginalBaseUnitPrice: r.BaseUnitPrice,
		OriginalPackPrice:     r.Batch.OriginalPackPrice,
		LineTotal:             r.LineTotal,
		RemainingBeforeExit:   r.RemainingBefore,
		RemainingAfterExit:    r.RemainingAfter,
		Sufficient:            r.RemainingAfter.GreaterThanOrEqual(decimal.Zero),
	}
}
