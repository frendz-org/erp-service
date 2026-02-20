package postgres

import (
	"context"

	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

func (r *baseRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := getTx(ctx); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}
