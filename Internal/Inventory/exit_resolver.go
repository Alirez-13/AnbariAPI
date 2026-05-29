package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ResolvedExitLine struct {
	Batch           *model.InventoryBatch
	Product         *model.Product
	Multiplier      decimal.Decimal
	InputQuantity   decimal.Decimal
	BaseQuantity    decimal.Decimal
	BaseUnitPrice   decimal.Decimal
	InputUnitPrice  decimal.Decimal
	LineTotal       decimal.Decimal
	RemainingBefore decimal.Decimal
	RemainingAfter  decimal.Decimal
}

type ExitResolver struct {
	repo Repository
}

func NewExitResolver(repo Repository) *ExitResolver {
	return &ExitResolver{repo: repo}
}

func (er *ExitResolver) Resolve(
	ctx context.Context,
	db *gorm.DB,
	lines []dto.ExitLineRequest,
	forUpdate bool,
) ([]ResolvedExitLine, error) {
	if len(lines) == 0 {
		return nil, ErrEmptyLines
	}

	resolved := make([]ResolvedExitLine, 0, len(lines))
	batchDeductions := make(map[uint]decimal.Decimal)
	productCache := make(map[uint]*model.Product)

	for i, line := range lines {
		batch, err := er.repo.GetBatch(ctx, db, line.BatchID, forUpdate)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i, err)
		}

		if batch.IsExhausted() {
			return nil, fmt.Errorf("%w: batch %d", ErrInsufficientStock, batch.ID)
		}

		product, ok := productCache[batch.ProductID]
		if !ok {
			product, err = er.repo.GetProduct(ctx, db, batch.ProductID)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", i, err)
			}
			productCache[batch.ProductID] = product
		}

		multiplier, err := resolveUnitMultiplier(ctx, db, er.repo, batch.ProductID, line.UnitName, product.BaseUnit)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i, err)
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

		resolved = append(resolved, ResolvedExitLine{
			Batch:           batch,
			Product:         product,
			Multiplier:      multiplier,
			InputQuantity:   line.Quantity,
			BaseQuantity:    baseQuantity,
			BaseUnitPrice:   baseUnitPrice,
			InputUnitPrice:  inputUnitPrice,
			LineTotal:       lineTotal,
			RemainingBefore: batch.RemainingBaseQuantity,
			RemainingAfter:  effectiveRemaining.Sub(baseQuantity),
		})
	}

	return resolved, nil
}

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