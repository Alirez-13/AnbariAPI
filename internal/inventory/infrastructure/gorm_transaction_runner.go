// # SINGLE REASON: Run inventory application work inside GORM transactions.
package infrastructure

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey struct{}

type GormTransactionRunner struct {
	db *gorm.DB
}

func NewGormTransactionRunner(db *gorm.DB) *GormTransactionRunner {
	return &GormTransactionRunner{db: db}
}

func (r *GormTransactionRunner) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ContextWithTx(ctx, tx))
	})
}

func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func DBFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback.WithContext(ctx)
}
