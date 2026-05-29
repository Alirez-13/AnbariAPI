package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"

	"github.com/shopspring/decimal"
)

func toBatchAvailabilityDTO(b model.InventoryBatch) dto.BatchAvailabilityDTO {
	return dto.BatchAvailabilityDTO{
		BatchID:               b.ID,
		EntryDate:             b.EntryDate,
		EntryUnitName:         b.EntryUnitName,
		EntryUnitMultiplier:   b.EntryUnitMultiplier,
		OriginalPackPrice:     b.OriginalPackPrice,
		OriginalBaseUnitPrice: b.OriginalBaseUnitPrice,
		RemainingBaseQuantity: b.RemainingBaseQuantity,
	}
}

func toTransactionDTO(txn *model.Transaction) *dto.TransactionDTO {
	details := make([]dto.TransactionDetailDTO, 0, len(txn.Details))
	for _, d := range txn.Details {
		productName := ""
		if d.Product.ID != 0 {
			productName = d.Product.Name
		}
		details = append(details, dto.TransactionDetailDTO{
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
	return &dto.TransactionDTO{
		ID:              txn.ID,
		TransactionType: txn.TransactionType,
		Reference:       txn.Reference,
		Date:            txn.Date,
		Details:         details,
		CreatedAt:       txn.CreatedAt,
	}
}

func toExitPreviewLineDTO(r resolvedExitLine, requestedUnit string) dto.ExitPreviewLineDTO {
	return dto.ExitPreviewLineDTO{
		BatchID:               r.batch.ID,
		ProductID:             r.product.ID,
		ProductName:           r.product.Name,
		EntryDate:             r.batch.EntryDate,
		RequestedUnit:         requestedUnit,
		RequestedQuantity:     r.inputQuantity,
		Multiplier:            r.multiplier,
		BaseQuantity:          r.baseQuantity,
		OriginalBaseUnitPrice: r.baseUnitPrice,
		OriginalPackPrice:     r.batch.OriginalPackPrice,
		LineTotal:             r.lineTotal,
		RemainingBeforeExit:   r.remainingBefore,
		RemainingAfterExit:    r.remainingAfter,
		Sufficient:            r.remainingAfter.GreaterThanOrEqual(decimal.Zero),
	}
}
