package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {

	if existingTx, ok := getTx(ctx); ok {

		return existingTx.Transaction(func(nestedTx *gorm.DB) error {
			nestedCtx := withTx(ctx, nestedTx)
			return fn(nestedCtx)
		})
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := withTx(ctx, tx)
		if err := fn(txCtx); err != nil {
			return fmt.Errorf("transaction failed: %w", err)
		}
		return nil
	})
}
