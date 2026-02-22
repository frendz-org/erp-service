package entity_test

import (
	"testing"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMasterdataItem_TableName(t *testing.T) {
	item := entity.MasterdataItem{}
	assert.Equal(t, "masterdata_items", item.TableName())
}

func TestMasterdataItem_Validate(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name    string
		item    *entity.MasterdataItem
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid item",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "MALE",
				Name:       "Male",
				Status:     entity.MasterdataItemStatusActive,
			},
			wantErr: false,
		},
		{
			name: "valid item with all fields",
			item: &entity.MasterdataItem{
				CategoryID:    categoryID,
				TenantID:      func() *uuid.UUID { id := uuid.New(); return &id }(),
				ParentItemID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
				Code:          "ID-JK",
				Name:          "DKI Jakarta",
				AltName:       func() *string { s := "Jakarta"; return &s }(),
				Description:   func() *string { s := "Special Capital Region"; return &s }(),
				SortOrder:     1,
				IsSystem:      true,
				IsDefault:     false,
				Status:        entity.MasterdataItemStatusActive,
				EffectiveFrom: func() *time.Time { t := time.Now().AddDate(-1, 0, 0); return &t }(),
				EffectiveUntil: func() *time.Time { t := time.Now().AddDate(1, 0, 0); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "missing category_id",
			item: &entity.MasterdataItem{
				Code:   "MALE",
				Name:   "Male",
				Status: entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "category_id is required",
		},
		{
			name: "missing code",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Name:       "Male",
				Status:     entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "code too long",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code: func() string {
					s := ""
					for i := 0; i < 101; i++ {
						s += "a"
					}
					return s
				}(),
				Name:   "Test",
				Status: entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "code must not exceed 100 characters",
		},
		{
			name: "missing name",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "MALE",
				Status:     entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "TEST",
				Name: func() string {
					s := ""
					for i := 0; i < 256; i++ {
						s += "a"
					}
					return s
				}(),
				Status: entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "name must not exceed 255 characters",
		},
		{
			name: "alt_name too long",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "TEST",
				Name:       "Test",
				AltName: func() *string {
					s := ""
					for i := 0; i < 256; i++ {
						s += "a"
					}
					return &s
				}(),
				Status: entity.MasterdataItemStatusActive,
			},
			wantErr: true,
			errMsg:  "alt_name must not exceed 255 characters",
		},
		{
			name: "invalid status",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "MALE",
				Name:       "Male",
				Status:     "INVALID",
			},
			wantErr: true,
			errMsg:  "status must be ACTIVE or INACTIVE",
		},
		{
			name: "self-referencing parent",
			item: func() *entity.MasterdataItem {
				id := uuid.New()
				return &entity.MasterdataItem{
					ID:           id,
					CategoryID:   categoryID,
					ParentItemID: &id,
					Code:         "SELF",
					Name:         "Self Reference",
					Status:       entity.MasterdataItemStatusActive,
				}
			}(),
			wantErr: true,
			errMsg:  "item cannot be its own parent",
		},
		{
			name: "invalid effective dates - from after until",
			item: &entity.MasterdataItem{
				CategoryID:     categoryID,
				Code:           "TEMP",
				Name:           "Temporary Item",
				Status:         entity.MasterdataItemStatusActive,
				EffectiveFrom:  func() *time.Time { t := time.Now().AddDate(1, 0, 0); return &t }(),
				EffectiveUntil: func() *time.Time { t := time.Now().AddDate(-1, 0, 0); return &t }(),
			},
			wantErr: true,
			errMsg:  "effective_from must be before or equal to effective_until",
		},
		{
			name: "valid inactive item",
			item: &entity.MasterdataItem{
				CategoryID: categoryID,
				Code:       "OLD",
				Name:       "Old Item",
				Status:     entity.MasterdataItemStatusInactive,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMasterdataItem_IsEffective(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	lastWeek := now.AddDate(0, 0, -7)
	nextWeek := now.AddDate(0, 0, 7)

	tests := []struct {
		name          string
		item          *entity.MasterdataItem
		checkTime     time.Time
		wantEffective bool
	}{
		{
			name: "active with no date restrictions",
			item: &entity.MasterdataItem{
				Status: entity.MasterdataItemStatusActive,
			},
			checkTime:     now,
			wantEffective: true,
		},
		{
			name: "inactive item",
			item: &entity.MasterdataItem{
				Status: entity.MasterdataItemStatusInactive,
			},
			checkTime:     now,
			wantEffective: false,
		},
		{
			name: "active within date range",
			item: &entity.MasterdataItem{
				Status:         entity.MasterdataItemStatusActive,
				EffectiveFrom:  &yesterday,
				EffectiveUntil: &tomorrow,
			},
			checkTime:     now,
			wantEffective: true,
		},
		{
			name: "before effective_from",
			item: &entity.MasterdataItem{
				Status:        entity.MasterdataItemStatusActive,
				EffectiveFrom: &tomorrow,
			},
			checkTime:     now,
			wantEffective: false,
		},
		{
			name: "after effective_until",
			item: &entity.MasterdataItem{
				Status:         entity.MasterdataItemStatusActive,
				EffectiveUntil: &yesterday,
			},
			checkTime:     now,
			wantEffective: false,
		},
		{
			name: "on effective_from date",
			item: &entity.MasterdataItem{
				Status:        entity.MasterdataItemStatusActive,
				EffectiveFrom: &lastWeek,
			},
			checkTime:     now,
			wantEffective: true,
		},
		{
			name: "on effective_until date",
			item: &entity.MasterdataItem{
				Status:         entity.MasterdataItemStatusActive,
				EffectiveUntil: &nextWeek,
			},
			checkTime:     now,
			wantEffective: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsEffective(tt.checkTime)
			assert.Equal(t, tt.wantEffective, result)
		})
	}
}

func TestMasterdataItem_IsCurrentlyEffective(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1)
	tomorrow := time.Now().AddDate(0, 0, 1)

	item := &entity.MasterdataItem{
		Status:         entity.MasterdataItemStatusActive,
		EffectiveFrom:  &yesterday,
		EffectiveUntil: &tomorrow,
	}

	assert.True(t, item.IsCurrentlyEffective())

	item.Status = entity.MasterdataItemStatusInactive
	assert.False(t, item.IsCurrentlyEffective())
}

func TestMasterdataItem_Deactivate(t *testing.T) {
	item := &entity.MasterdataItem{
		CategoryID: uuid.New(),
		Code:       "MALE",
		Name:       "Male",
		Status:     entity.MasterdataItemStatusActive,
	}

	item.Deactivate()

	assert.Equal(t, entity.MasterdataItemStatusInactive, item.Status)
	assert.False(t, item.UpdatedAt.IsZero())
}

func TestMasterdataItem_Activate(t *testing.T) {
	item := &entity.MasterdataItem{
		CategoryID: uuid.New(),
		Code:       "MALE",
		Name:       "Male",
		Status:     entity.MasterdataItemStatusInactive,
	}

	item.Activate()

	assert.Equal(t, entity.MasterdataItemStatusActive, item.Status)
	assert.False(t, item.UpdatedAt.IsZero())
}

func TestMasterdataItem_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status entity.MasterdataItemStatus
		want   bool
	}{
		{"active", entity.MasterdataItemStatusActive, true},
		{"inactive", entity.MasterdataItemStatusInactive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &entity.MasterdataItem{Status: tt.status}
			assert.Equal(t, tt.want, item.IsActive())
		})
	}
}

func TestMasterdataItem_IsGlobal(t *testing.T) {
	tests := []struct {
		name     string
		tenantID *uuid.UUID
		want     bool
	}{
		{
			name:     "global item (no tenant)",
			tenantID: nil,
			want:     true,
		},
		{
			name:     "tenant-specific item",
			tenantID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &entity.MasterdataItem{TenantID: tt.tenantID}
			assert.Equal(t, tt.want, item.IsGlobal())
			assert.Equal(t, !tt.want, item.IsTenantSpecific())
		})
	}
}

func TestMasterdataItem_HasParent(t *testing.T) {
	tests := []struct {
		name         string
		parentItemID *uuid.UUID
		want         bool
	}{
		{
			name:         "with parent",
			parentItemID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			want:         true,
		},
		{
			name:         "without parent (root item)",
			parentItemID: nil,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &entity.MasterdataItem{ParentItemID: tt.parentItemID}
			assert.Equal(t, tt.want, item.HasParent())
			assert.Equal(t, !tt.want, item.IsRootItem())
		})
	}
}

func TestMasterdataItem_IncrementVersion(t *testing.T) {
	item := &entity.MasterdataItem{
		CategoryID: uuid.New(),
		Code:       "MALE",
		Name:       "Male",
		Version:    1,
	}

	originalVersion := item.Version
	item.IncrementVersion()

	assert.Equal(t, originalVersion+1, item.Version)
	assert.False(t, item.UpdatedAt.IsZero())
}

func TestMasterdataItem_StatusConstants(t *testing.T) {
	assert.Equal(t, entity.MasterdataItemStatus("ACTIVE"), entity.MasterdataItemStatusActive)
	assert.Equal(t, entity.MasterdataItemStatus("INACTIVE"), entity.MasterdataItemStatusInactive)
}
