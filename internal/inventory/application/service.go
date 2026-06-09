package application

// # SINGLE REASON: Define inventory application service contract and dependencies.

import (
	"context"

	models "AnbariAPI/internal/inventory/domain"
)

type InventoryService interface {
	ProcessInboundTransaction(ctx context.Context, input InboundDTO) (*models.Transaction, error)
	ProcessOutboundTransaction(ctx context.Context, input OutboundDTO) (*models.Transaction, error)
}

type InventoryServiceImpl struct {
	repo    models.InventoryRepository
	tx      models.TransactionRunner
	retries int
}

func NewInventoryService(repo models.InventoryRepository, tx models.TransactionRunner) *InventoryServiceImpl {
	return &InventoryServiceImpl{repo: repo, tx: tx, retries: 3}
}
