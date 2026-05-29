package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"

	"github.com/shopspring/decimal"
)

func ToBatchAvailabilityDTO(b model.InventoryBatch) dto.BatchAvailabilityDTO {
	return dto.BatchAvailabilityDTO{
		BatchID:               b.ID,
		EntryDate:             b.EntryDate,
		EntryUnitName:         b.EntryUnitName,
		EntryUnitMultiplier:     b.EntryUnitMultiplier,
		OriginalPackPrice:       b.OriginalPackPrice,
		OriginalBaseUnitPrice:   b.OriginalBaseUnitPrice,
		RemainingBaseQuantity:   b.RemainingBaseQuantity,
	}
}

func ToTransactionDTO(txn *model.Transaction) *dto.TransactionDTO {
	if txn == nil {
		return nil
	}
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

func ToExitPreviewLineDTO(r ResolvedExitLine, requestedUnit string) dto.ExitPreviewLineDTO {
	return dto.ExitPreviewLineDTO{
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