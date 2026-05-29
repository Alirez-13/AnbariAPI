package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type resolvedExitLine struct {
	batch           *model.InventoryBatch
	product         *model.Product
	multiplier      decimal.Decimal
	inputQuantity   decimal.Decimal
	baseQuantity    decimal.Decimal
	baseUnitPrice   decimal.Decimal
	inputUnitPrice  decimal.Decimal
	lineTotal       decimal.Decimal
	remainingBefore decimal.Decimal
	remainingAfter  decimal.Decimal
}

type exitResolver struct {
	repo Repository
}

func newExitResolver(repo Repository) *exitResolver {
	return &exitResolver{repo: repo}
}

func (er *exitResolver) resolve(
	ctx context.Context,
	db *gorm.DB,
	lines []dto.ExitLineRequest,
	forUpdate bool,
) ([]resolvedExitLine, error) {
	resolved := make([]resolvedExitLine, 0, len(lines))
	batchDeductions := make(map[uint]decimal.Decimal)
	productCache := make(map[uint]*model.Product)

	for _, line := range lines {
		batch, err := er.repo.GetBatch(ctx, db, line.BatchID, forUpdate)
		if err != nil {
			return nil, err
		}

		if batch.RemainingBaseQuantity.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("%w: batch %d is exhausted", ErrInsufficientStock, batch.ID)
		}

		product, ok := productCache[batch.ProductID]
		if !ok {
			product, err = er.repo.GetProduct(ctx, db, batch.ProductID)
			if err != nil {
				return nil, err
			}
			productCache[batch.ProductID] = product
		}

		multiplier, err := resolveUnitMultiplier(ctx, db, er.repo, batch.ProductID, line.UnitName, product.BaseUnit)
		if err != nil {
			return nil, err
		}

		baseQuantity := line.Quantity.Mul(multiplier)
		baseUnitPrice := batch.OriginalBaseUnitPrice
		inputUnitPrice := baseUnitPrice.Mul(multiplier)
		lineTotal := baseQuantity.Mul(baseUnitPrice)

		prevDeduction := batchDeductions[batch.ID]
		effectiveRemaining := batch.RemainingBaseQuantity.Sub(prevDeduction)

		if effectiveRemaining.LessThan(baseQuantity) {
			return nil, fmt.Errorf(
				"%w: batch %d has %s base units available, but %s requested",
				ErrInsufficientStock,
				batch.ID,
				effectiveRemaining.StringFixed(4),
				baseQuantity.StringFixed(4),
			)
		}

		batchDeductions[batch.ID] = prevDeduction.Add(baseQuantity)

		resolved = append(resolved, resolvedExitLine{
			batch:           batch,
			product:         product,
			multiplier:      multiplier,
			inputQuantity:   line.Quantity,
			baseQuantity:    baseQuantity,
			baseUnitPrice:   baseUnitPrice,
			inputUnitPrice:  inputUnitPrice,
			lineTotal:       lineTotal,
			remainingBefore: batch.RemainingBaseQuantity,
			remainingAfter:  effectiveRemaining.Sub(baseQuantity),
		})
	}

	return resolved, nil
}

// resolveUnitMultiplier is a pure calculation helper — no receiver needed.
func resolveUnitMultiplier(
	ctx context.Context,
	db *gorm.DB,
	repo Repository,
	productID uint,
	unitName string,
	baseUnit string,
) (decimal.Decimal, error) {
	if unitName == baseUnit {
		return decimal.NewFromInt(1), nil
	}
	pu, err := repo.GetProductUnit(ctx, db, productID, unitName)
	if err != nil {
		return decimal.Zero, err
	}
	return pu.Multiplier, nil
}
