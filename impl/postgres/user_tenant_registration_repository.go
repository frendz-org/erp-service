package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"iam-service/entity"
	apperrors "iam-service/pkg/errors"
	membercontract "iam-service/saving/member/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userTenantRegistrationRepository struct {
	baseRepository
}

func NewUserTenantRegistrationRepository(db *gorm.DB) *userTenantRegistrationRepository {
	return &userTenantRegistrationRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userTenantRegistrationRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserTenantRegistration, error) {
	var registrations []entity.UserTenantRegistration
	err := r.getDB(ctx).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, entity.UTRStatusActive).
		Find(&registrations).Error
	if err != nil {
		return nil, translateError(err, "user tenant registration")
	}
	return registrations, nil
}

func (r *userTenantRegistrationRepository) Create(ctx context.Context, reg *entity.UserTenantRegistration) error {
	if err := r.getDB(ctx).Create(reg).Error; err != nil {
		return translateError(err, "member registration")
	}
	return nil
}

func (r *userTenantRegistrationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UserTenantRegistration, error) {
	var reg entity.UserTenantRegistration
	err := r.getDB(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&reg).Error
	if err != nil {
		return nil, translateError(err, "member registration")
	}
	return &reg, nil
}

func (r *userTenantRegistrationRepository) GetByUserAndProduct(ctx context.Context, userID, tenantID, productID uuid.UUID, regType string) (*entity.UserTenantRegistration, error) {
	var reg entity.UserTenantRegistration
	err := r.getDB(ctx).
		Where("user_id = ? AND tenant_id = ? AND product_id = ? AND registration_type = ? AND deleted_at IS NULL",
			userID, tenantID, productID, regType).
		First(&reg).Error
	if err != nil {
		return nil, translateError(err, "member registration")
	}
	return &reg, nil
}

func (r *userTenantRegistrationRepository) UpdateStatus(ctx context.Context, reg *entity.UserTenantRegistration) error {
	result := r.getDB(ctx).
		Model(&entity.UserTenantRegistration{}).
		Where("id = ? AND version = ? AND deleted_at IS NULL", reg.ID, reg.Version).
		Updates(map[string]interface{}{
			"status":      reg.Status,
			"approved_by": reg.ApprovedBy,
			"approved_at": reg.ApprovedAt,
			"metadata":    reg.Metadata,
			"version":     gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return translateError(result.Error, "member registration")
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrConflict("member registration was modified by another request, please retry")
	}

	reg.Version++
	return nil
}

var allowedMemberSortColumns = map[string]string{
	"created_at": "utr.created_at",
	"status":     "utr.status",
	"first_name": "up.first_name",
	"last_name":  "up.last_name",
}

func (r *userTenantRegistrationRepository) ListByProductWithFilters(ctx context.Context, filter *membercontract.MemberListFilter) ([]membercontract.MemberListRow, int64, error) {
	db := r.getDB(ctx)

	baseQuery := db.Table("user_tenant_registrations AS utr").
		Joins("JOIN users u ON u.id = utr.user_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN user_profiles up ON up.user_id = utr.user_id").
		Joins("LEFT JOIN user_role_assignments ura ON ura.user_id = utr.user_id AND ura.product_id = utr.product_id AND ura.deleted_at IS NULL").
		Joins("LEFT JOIN roles r ON r.id = ura.role_id AND r.deleted_at IS NULL AND r.status = 'ACTIVE'").
		Where("utr.tenant_id = ? AND utr.product_id = ? AND utr.registration_type = ? AND utr.deleted_at IS NULL",
			filter.TenantID, filter.ProductID, "MEMBER")

	if filter.Status != nil && *filter.Status != "" {
		baseQuery = baseQuery.Where("utr.status = ?", *filter.Status)
	}

	if filter.Search != "" {
		search := "%" + escapeILIKE(filter.Search) + "%"
		baseQuery = baseQuery.Where("(up.first_name ILIKE ? OR up.last_name ILIKE ? OR u.email ILIKE ?)",
			search, search, search)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "member registration")
	}

	sortCol := "utr.created_at"
	if col, ok := allowedMemberSortColumns[filter.SortBy]; ok {
		sortCol = col
	}
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}
	orderClause := fmt.Sprintf("%s %s", sortCol, sortOrder)

	offset := (filter.Page - 1) * filter.PerPage

	type scanRow struct {
		ID               uuid.UUID  `gorm:"column:id"`
		UserID           uuid.UUID  `gorm:"column:user_id"`
		TenantID         uuid.UUID  `gorm:"column:tenant_id"`
		ProductID        *uuid.UUID `gorm:"column:product_id"`
		RegistrationType string     `gorm:"column:registration_type"`
		Status           string     `gorm:"column:utr_status"`
		CreatedAt        time.Time  `gorm:"column:utr_created_at"`
		FirstName        string     `gorm:"column:first_name"`
		LastName         string     `gorm:"column:last_name"`
		Email            string     `gorm:"column:email"`
		RoleCode         *string    `gorm:"column:role_code"`
		RoleName         *string    `gorm:"column:role_name"`
	}

	var scanRows []scanRow
	err := baseQuery.
		Select("utr.id, utr.user_id, utr.tenant_id, utr.product_id, utr.registration_type, utr.status AS utr_status, utr.created_at AS utr_created_at, COALESCE(up.first_name, '') AS first_name, COALESCE(up.last_name, '') AS last_name, u.email, r.code AS role_code, r.name AS role_name").
		Order(orderClause).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&scanRows).Error
	if err != nil {
		return nil, 0, translateError(err, "member registration")
	}

	rows := make([]membercontract.MemberListRow, 0, len(scanRows))
	for _, sr := range scanRows {
		row := membercontract.MemberListRow{
			Registration: entity.UserTenantRegistration{
				ID:               sr.ID,
				UserID:           sr.UserID,
				TenantID:         sr.TenantID,
				ProductID:        sr.ProductID,
				RegistrationType: sr.RegistrationType,
				Status:           entity.UserTenantRegistrationStatus(sr.Status),
				CreatedAt:        sr.CreatedAt,
			},
			FirstName: sr.FirstName,
			LastName:  sr.LastName,
			Email:     sr.Email,
			RoleCode:  sr.RoleCode,
			RoleName:  sr.RoleName,
		}
		rows = append(rows, row)
	}

	return rows, total, nil
}
