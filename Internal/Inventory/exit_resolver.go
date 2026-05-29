package Inventory

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"context"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
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

// ExitResolver orchestrates the domain logic for mapping exit requests to actual inventory records.
type ExitResolver interface {
	Resolve(ctx context.Context, repo Repository, lines []dto.ExitLineRequest, forUpdate bool) ([]ResolvedExitLine, error)
}

type exitResolverImpl struct{}

// NewExitResolver creates a new instance of the ExitResolver domain service.
func NewExitResolver() ExitResolver {
	return &exitResolverImpl{}
}

func (er *exitResolverImpl) Resolve(
	ctx context.Context,
	repo Repository,
	lines []dto.ExitLineRequest,
	forUpdate bool,
) ([]ResolvedExitLine, error) {
	if len(lines) == 0 {
		return nil, ErrEmptyLines
	}

	// 1. Prevent Deadlocks: Collect unique batch IDs and sort them to ensure deterministic locking order.
	batchIDMap := make(map[uint]struct{})
	for _, line := range lines {
		batchIDMap[line.BatchID] = struct{}{}
	}

	sortedBatchIDs := make([]uint, 0, len(batchIDMap))
	for id := range batchIDMap {
		sortedBatchIDs = append(sortedBatchIDs, id)
	}
	sort.Slice(sortedBatchIDs, func(i, j int) bool { return sortedBatchIDs[i] < sortedBatchIDs[j] })

	// 2. Pre-fetch and lock batches safely
	batchCache := make(map[uint]*model.InventoryBatch, len(sortedBatchIDs))
	for _, id := range sortedBatchIDs {
		var batch *model.InventoryBatch
		var err error

		if forUpdate {
			batch, err = repo.GetBatchForUpdate(ctx, id)
		} else {
			batch, err = repo.GetBatch(ctx, id)
		}

		if err != nil {
			return nil, fmt.Errorf("retrieve batch %d: %w", id, err)
		}
		if batch.IsExhausted() {
			return nil, fmt.Errorf("%w: batch %d is exhausted", ErrInsufficientStock, batch.ID)
		}
		batchCache[id] = batch
	}

	resolved := make([]ResolvedExitLine, 0, len(lines))
	batchDeductions := make(map[uint]decimal.Decimal)

	// Local caches to prevent N+1 queries during resolution loop
	productCache := make(map[uint]*model.Product)
	unitCache := make(map[string]decimal.Decimal) // Key format: "productID:unitName"

	// 3. Process lines in original order to maintain 1:1 mapping with request
	for i, line := range lines {
		batch := batchCache[line.BatchID]

		product, ok := productCache[batch.ProductID]
		if !ok {
			p, err := repo.GetProduct(ctx, batch.ProductID)
			if err != nil {
				return nil, fmt.Errorf("line %d product lookup: %w", i, err)
			}
			product = p
			productCache[batch.ProductID] = product
		}

		// Resolve Multiplier with local cache
		cacheKey := fmt.Sprintf("%d:%s", product.ID, line.UnitName)
		multiplier, cached := unitCache[cacheKey]
		if !cached {
			m, err := resolveUnitMultiplier(ctx, repo, product.ID, line.UnitName, product.BaseUnit)
			if err != nil {
				return nil, fmt.Errorf("line %d multiplier lookup: %w", i, err)
			}
			multiplier = m
			unitCache[cacheKey] = multiplier
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
	repo Repository,
	productID uint,
	unitName string,
	baseUnit string,
) (decimal.Decimal, error) {
	if unitName == baseUnit {
		return decimal.NewFromInt(1), nil
	}
	pu, err := repo.GetProductUnit(ctx, productID, unitName)
	if err != nil {
		return decimal.Zero, err
	}
	return pu.Multiplier, nil
}
