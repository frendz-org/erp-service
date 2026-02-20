package internal

import (
	"context"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	item, err := uc.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrNotFound("item not found")
		}
		return err
	}

	if item.IsSystem {
		return errors.ErrValidation("system items cannot be deleted")
	}

	if err := uc.itemRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = uc.cache.InvalidateItems(ctx)

	return nil
}
