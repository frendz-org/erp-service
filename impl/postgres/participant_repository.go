package postgres

import (
	"context"
	"fmt"
	"strings"

	"erp-service/entity"
	apperrors "erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var allowedParticipantSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"full_name":  true,
	"status":     true,
}

var ilikeReplacer = strings.NewReplacer(
	"\\", "\\\\",
	"%", "\\%",
	"_", "\\_",
)

func escapeILIKE(s string) string {
	return ilikeReplacer.Replace(s)
}

type participantRepository struct {
	baseRepository
}

func NewParticipantRepository(db *gorm.DB) participant.ParticipantRepository {
	return &participantRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantRepository) Create(ctx context.Context, participant *entity.Participant) error {
	if err := r.getDB(ctx).Create(participant).Error; err != nil {
		return translateError(err, "participant")
	}
	return nil
}
func (r *participantRepository) GetByKTPAndPensionNumber(ctx context.Context, ktpNumber, pensionNumber string, tenantID, productID uuid.UUID) (*entity.Participant, *entity.ParticipantPension, error) {
	type joinedRow struct {
		entity.Participant   `gorm:"embedded"`
		PensionID            uuid.UUID `gorm:"column:pension_id"`
		PensionParticipantID uuid.UUID `gorm:"column:pension_participant_id"`
		PensionNumber        *string   `gorm:"column:pension_participant_number"`
	}

	var row joinedRow
	err := r.getDB(ctx).Raw(`
		SELECT
			p.*,
			pp.id                  AS pension_id,
			pp.participant_id      AS pension_participant_id,
			pp.participant_number  AS pension_participant_number
		FROM participants p
		JOIN participant_pensions pp
		  ON pp.participant_id = p.id
		  AND pp.deleted_at IS NULL
		WHERE p.ktp_number     = ?
		  AND pp.participant_number = ?
		  AND p.tenant_id      = ?
		  AND p.product_id     = ?
		  AND p.deleted_at     IS NULL
		LIMIT 1
	`, ktpNumber, pensionNumber, tenantID, productID).Scan(&row).Error

	if err != nil {
		return nil, nil, fmt.Errorf("GetByKTPAndPensionNumber: %w", err)
	}

	if row.Participant.ID == (uuid.UUID{}) {
		return nil, nil, apperrors.ErrNotFound("participant not found")
	}

	pension := &entity.ParticipantPension{
		ID:                row.PensionID,
		ParticipantID:     row.PensionParticipantID,
		ParticipantNumber: row.PensionNumber,
	}
	return &row.Participant, pension, nil
}

func (r *participantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&participant).Error
	if err != nil {
		return nil, translateError(err, "participant")
	}
	return &participant, nil
}

func (r *participantRepository) GetByKTPNumber(ctx context.Context, tenantID, productID uuid.UUID, ktpNumber string) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.getDB(ctx).
		Where("tenant_id = ? AND product_id = ? AND ktp_number = ? AND deleted_at IS NULL", tenantID, productID, ktpNumber).
		First(&participant).Error
	if err != nil {
		return nil, translateError(err, "participant")
	}
	return &participant, nil
}

func (r *participantRepository) GetByEmployeeNumber(ctx context.Context, tenantID, productID uuid.UUID, employeeNumber string) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.getDB(ctx).
		Where("tenant_id = ? AND product_id = ? AND employee_number = ? AND deleted_at IS NULL", tenantID, productID, employeeNumber).
		First(&participant).Error
	if err != nil {
		return nil, translateError(err, "participant")
	}
	return &participant, nil
}

func (r *participantRepository) GetByUserAndTenantProduct(ctx context.Context, userID, tenantID, productID uuid.UUID) (*entity.Participant, error) {
	var p entity.Participant
	err := r.getDB(ctx).
		Where("user_id = ? AND tenant_id = ? AND product_id = ? AND deleted_at IS NULL", userID, tenantID, productID).
		First(&p).Error
	if err != nil {
		return nil, translateError(err, "participant")
	}
	return &p, nil
}

func (r *participantRepository) Update(ctx context.Context, participant *entity.Participant) error {
	oldVersion := participant.Version
	participant.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(participant)
	if result.Error != nil {
		participant.Version = oldVersion
		return translateError(result.Error, "participant")
	}
	if result.RowsAffected == 0 {
		participant.Version = oldVersion
		return apperrors.ErrConflict("participant was modified by another request")
	}
	return nil
}

func (r *participantRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.Participant{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant")
	}
	return nil
}

func (r *participantRepository) List(ctx context.Context, filter *participant.ParticipantFilter) ([]*entity.Participant, int64, error) {
	var participants []*entity.Participant
	var total int64

	query := r.getDB(ctx).Model(&entity.Participant{}).
		Where("tenant_id = ? AND product_id = ? AND deleted_at IS NULL", filter.TenantID, filter.ProductID)

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != "" {
		search := "%" + escapeILIKE(filter.Search) + "%"
		query = query.Where(
			"full_name ILIKE ? OR ktp_number ILIKE ? OR employee_number ILIKE ? OR phone_number ILIKE ?",
			search, search, search, search,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "participant")
	}

	sortBy := "created_at"
	if filter.SortBy != "" && allowedParticipantSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}
	sortOrder := "desc"
	if filter.SortOrder == "asc" {
		sortOrder = "asc"
	}

	offset := (filter.Page - 1) * filter.PerPage
	err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&participants).Error

	if err != nil {
		return nil, 0, translateError(err, "participant")
	}

	return participants, total, nil
}
