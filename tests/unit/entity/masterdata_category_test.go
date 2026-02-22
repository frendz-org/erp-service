package entity_test

import (
	"testing"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMasterdataCategory_TableName(t *testing.T) {
	cat := entity.MasterdataCategory{}
	assert.Equal(t, "masterdata_categories", cat.TableName())
}

func TestMasterdataCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cat     *entity.MasterdataCategory
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flat category",
			cat: &entity.MasterdataCategory{
				Code:   "GENDER",
				Name:   "Gender",
				Status: entity.MasterdataCategoryStatusActive,
			},
			wantErr: false,
		},
		{
			name: "valid hierarchical category",
			cat: &entity.MasterdataCategory{
				ID:               uuid.New(),
				Code:             "PROVINCE",
				Name:             "Province",
				ParentCategoryID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				Status:           entity.MasterdataCategoryStatusActive,
			},
			wantErr: false,
		},
		{
			name: "missing code",
			cat: &entity.MasterdataCategory{
				Name:   "Gender",
				Status: entity.MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "code too long",
			cat: &entity.MasterdataCategory{
				Code:   "THIS_CODE_IS_WAY_TOO_LONG_AND_EXCEEDS_THE_FIFTY_CHARACTER_LIMIT_SET_BY_DATABASE",
				Name:   "Test",
				Status: entity.MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "code must not exceed 50 characters",
		},
		{
			name: "missing name",
			cat: &entity.MasterdataCategory{
				Code:   "GENDER",
				Status: entity.MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			cat: &entity.MasterdataCategory{
				Code: "TEST",
				Name: func() string {
					s := ""
					for i := 0; i < 256; i++ {
						s += "a"
					}
					return s
				}(),
				Status: entity.MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "name must not exceed 255 characters",
		},
		{
			name: "invalid status",
			cat: &entity.MasterdataCategory{
				Code:   "GENDER",
				Name:   "Gender",
				Status: "INVALID",
			},
			wantErr: true,
			errMsg:  "status must be ACTIVE or INACTIVE",
		},
		{
			name: "self-referencing parent",
			cat: func() *entity.MasterdataCategory {
				id := uuid.New()
				return &entity.MasterdataCategory{
					ID:               id,
					Code:             "SELF",
					Name:             "Self Reference",
					ParentCategoryID: &id,
					Status:           entity.MasterdataCategoryStatusActive,
				}
			}(),
			wantErr: true,
			errMsg:  "category cannot be its own parent",
		},
		{
			name: "valid inactive category",
			cat: &entity.MasterdataCategory{
				Code:   "OLD_CATEGORY",
				Name:   "Old Category",
				Status: entity.MasterdataCategoryStatusInactive,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cat.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMasterdataCategory_Deactivate(t *testing.T) {
	cat := &entity.MasterdataCategory{
		Code:   "GENDER",
		Name:   "Gender",
		Status: entity.MasterdataCategoryStatusActive,
	}

	cat.Deactivate()

	assert.Equal(t, entity.MasterdataCategoryStatusInactive, cat.Status)
	assert.False(t, cat.UpdatedAt.IsZero())
}

func TestMasterdataCategory_Activate(t *testing.T) {
	cat := &entity.MasterdataCategory{
		Code:   "GENDER",
		Name:   "Gender",
		Status: entity.MasterdataCategoryStatusInactive,
	}

	cat.Activate()

	assert.Equal(t, entity.MasterdataCategoryStatusActive, cat.Status)
	assert.False(t, cat.UpdatedAt.IsZero())
}

func TestMasterdataCategory_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status entity.MasterdataCategoryStatus
		want   bool
	}{
		{"active", entity.MasterdataCategoryStatusActive, true},
		{"inactive", entity.MasterdataCategoryStatusInactive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := &entity.MasterdataCategory{Status: tt.status}
			assert.Equal(t, tt.want, cat.IsActive())
		})
	}
}

func TestMasterdataCategory_IsHierarchical(t *testing.T) {
	tests := []struct {
		name             string
		parentCategoryID *uuid.UUID
		want             bool
	}{
		{
			name:             "with parent",
			parentCategoryID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			want:             true,
		},
		{
			name:             "without parent",
			parentCategoryID: nil,
			want:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := &entity.MasterdataCategory{ParentCategoryID: tt.parentCategoryID}
			assert.Equal(t, tt.want, cat.IsHierarchical())
			assert.Equal(t, !tt.want, cat.IsFlat())
		})
	}
}

func TestMasterdataCategory_IncrementVersion(t *testing.T) {
	cat := &entity.MasterdataCategory{
		Code:    "GENDER",
		Name:    "Gender",
		Version: 1,
	}

	originalVersion := cat.Version
	cat.IncrementVersion()

	assert.Equal(t, originalVersion+1, cat.Version)
	assert.False(t, cat.UpdatedAt.IsZero())
}

func TestMasterdataCategory_StatusConstants(t *testing.T) {
	assert.Equal(t, entity.MasterdataCategoryStatus("ACTIVE"), entity.MasterdataCategoryStatusActive)
	assert.Equal(t, entity.MasterdataCategoryStatus("INACTIVE"), entity.MasterdataCategoryStatusInactive)
}

func TestMasterdataCategory_JSONTags(t *testing.T) {
	now := time.Now()
	description := "Test description"
	parentID := uuid.New()

	cat := entity.MasterdataCategory{
		ID:                 uuid.New(),
		Code:               "TEST",
		Name:               "Test Category",
		Description:        &description,
		ParentCategoryID:   &parentID,
		IsSystem:           true,
		IsTenantExtensible: false,
		SortOrder:          1,
		Status:             entity.MasterdataCategoryStatusActive,
		Version:            1,
		Timestamps: entity.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	assert.NotEmpty(t, cat.ID)
	assert.Equal(t, "TEST", cat.Code)
	assert.Equal(t, "Test Category", cat.Name)
	assert.Equal(t, &description, cat.Description)
	assert.Equal(t, &parentID, cat.ParentCategoryID)
	assert.True(t, cat.IsSystem)
	assert.False(t, cat.IsTenantExtensible)
}
