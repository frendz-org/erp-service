package internal

import (
	"context"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrNotFound("category not found")
		}
		return err
	}

	if category.IsSystem {
		return errors.ErrValidation("system categories cannot be deleted")
	}

	if err := uc.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = uc.cache.InvalidateCategories(ctx)

	return nil
}
