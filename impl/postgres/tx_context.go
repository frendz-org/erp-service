package postgres

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func getTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

func hasTx(ctx context.Context) bool {
	_, ok := getTx(ctx)
	return ok
}
