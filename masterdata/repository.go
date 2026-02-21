package masterdata

import (
	"context"

	"erp-service/entity"

	"github.com/google/uuid"
)

type CategoryFilter struct {
	Status    string
	IsSystem  *bool
	ParentID  *uuid.UUID
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

type ItemFilter struct {
	CategoryID   *uuid.UUID
	CategoryCode string
	TenantID     *uuid.UUID
	ParentID     *uuid.UUID
	ParentCode   string
	Status       string
	Search       string
	IsDefault    *bool
	IsSystem     *bool
	Page         int
	PerPage      int
	SortBy       string
	SortOrder    string
}

type CategoryRepository interface {
	List(ctx context.Context, filter *CategoryFilter) ([]*entity.MasterdataCategory, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataCategory, error)
	GetByCode(ctx context.Context, code string) (*entity.MasterdataCategory, error)

	Create(ctx context.Context, category *entity.MasterdataCategory) error
	Update(ctx context.Context, category *entity.MasterdataCategory) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataCategory, error)

	ExistsByCode(ctx context.Context, code string) (bool, error)
}

type ItemRepository interface {
	List(ctx context.Context, filter *ItemFilter) ([]*entity.MasterdataItem, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataItem, error)
	GetByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*entity.MasterdataItem, error)

	ValidateCode(ctx context.Context, categoryCode string, itemCode string, tenantID *uuid.UUID) (bool, error)

	Create(ctx context.Context, item *entity.MasterdataItem) error
	Update(ctx context.Context, item *entity.MasterdataItem) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataItem, error)
	GetTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error)
	ListByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error)

	ExistsByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (bool, error)
	GetDefaultItem(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*entity.MasterdataItem, error)
}
