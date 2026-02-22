package member

import (
	"context"

	"erp-service/entity"

	"github.com/google/uuid"
)

type MemberListFilter struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Status    *string
	Search    string
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

type MemberListRow struct {
	Registration entity.UserTenantRegistration
	FirstName    string
	LastName     string
	Email        string
	RoleCode     *string
	RoleName     *string
}

type UserTenantRegistrationRepository interface {
	Create(ctx context.Context, reg *entity.UserTenantRegistration) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UserTenantRegistration, error)
	GetByUserAndProduct(ctx context.Context, userID, tenantID, productID uuid.UUID, regType string) (*entity.UserTenantRegistration, error)
	UpdateStatus(ctx context.Context, reg *entity.UserTenantRegistration) error
	ListByProductWithFilters(ctx context.Context, filter *MemberListFilter) ([]MemberListRow, int64, error)
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole *entity.UserRole) error
	SoftDeleteByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) error
	GetActiveByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*entity.UserRole, error)
}

type ProductRepository interface {
	GetByIDAndTenant(ctx context.Context, productID, tenantID uuid.UUID) (*entity.Product, error)
}

type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetByCodeAndProduct(ctx context.Context, productID uuid.UUID, code string) (*entity.Role, error)
}

type ProductRegistrationConfigRepository interface {
	GetByProductAndType(ctx context.Context, productID uuid.UUID, regType string) (*entity.ProductRegistrationConfig, error)
}

type UserProfileRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
