package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMasterdataCategory_TableName(t *testing.T) {
	cat := MasterdataCategory{}
	assert.Equal(t, "masterdata_categories", cat.TableName())
}

func TestMasterdataCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cat     *MasterdataCategory
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flat category",
			cat: &MasterdataCategory{
				Code:   "GENDER",
				Name:   "Gender",
				Status: MasterdataCategoryStatusActive,
			},
			wantErr: false,
		},
		{
			name: "valid hierarchical category",
			cat: &MasterdataCategory{
				ID:               uuid.New(),
				Code:             "PROVINCE",
				Name:             "Province",
				ParentCategoryID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				Status:           MasterdataCategoryStatusActive,
			},
			wantErr: false,
		},
		{
			name: "missing code",
			cat: &MasterdataCategory{
				Name:   "Gender",
				Status: MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "code too long",
			cat: &MasterdataCategory{
				Code:   "THIS_CODE_IS_WAY_TOO_LONG_AND_EXCEEDS_THE_FIFTY_CHARACTER_LIMIT_SET_BY_DATABASE",
				Name:   "Test",
				Status: MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "code must not exceed 50 characters",
		},
		{
			name: "missing name",
			cat: &MasterdataCategory{
				Code:   "GENDER",
				Status: MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			cat: &MasterdataCategory{
				Code: "TEST",
				Name: func() string {
					s := ""
					for i := 0; i < 256; i++ {
						s += "a"
					}
					return s
				}(),
				Status: MasterdataCategoryStatusActive,
			},
			wantErr: true,
			errMsg:  "name must not exceed 255 characters",
		},
		{
			name: "invalid status",
			cat: &MasterdataCategory{
				Code:   "GENDER",
				Name:   "Gender",
				Status: "INVALID",
			},
			wantErr: true,
			errMsg:  "status must be ACTIVE or INACTIVE",
		},
		{
			name: "self-referencing parent",
			cat: func() *MasterdataCategory {
				id := uuid.New()
				return &MasterdataCategory{
					ID:               id,
					Code:             "SELF",
					Name:             "Self Reference",
					ParentCategoryID: &id,
					Status:           MasterdataCategoryStatusActive,
				}
			}(),
			wantErr: true,
			errMsg:  "category cannot be its own parent",
		},
		{
			name: "valid inactive category",
			cat: &MasterdataCategory{
				Code:   "OLD_CATEGORY",
				Name:   "Old Category",
				Status: MasterdataCategoryStatusInactive,
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
	cat := &MasterdataCategory{
		Code:   "GENDER",
		Name:   "Gender",
		Status: MasterdataCategoryStatusActive,
	}

	cat.Deactivate()

	assert.Equal(t, MasterdataCategoryStatusInactive, cat.Status)
	assert.False(t, cat.UpdatedAt.IsZero())
}

func TestMasterdataCategory_Activate(t *testing.T) {
	cat := &MasterdataCategory{
		Code:   "GENDER",
		Name:   "Gender",
		Status: MasterdataCategoryStatusInactive,
	}

	cat.Activate()

	assert.Equal(t, MasterdataCategoryStatusActive, cat.Status)
	assert.False(t, cat.UpdatedAt.IsZero())
}

func TestMasterdataCategory_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status MasterdataCategoryStatus
		want   bool
	}{
		{"active", MasterdataCategoryStatusActive, true},
		{"inactive", MasterdataCategoryStatusInactive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := &MasterdataCategory{Status: tt.status}
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
			cat := &MasterdataCategory{ParentCategoryID: tt.parentCategoryID}
			assert.Equal(t, tt.want, cat.IsHierarchical())
			assert.Equal(t, !tt.want, cat.IsFlat())
		})
	}
}

func TestMasterdataCategory_IncrementVersion(t *testing.T) {
	cat := &MasterdataCategory{
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
	assert.Equal(t, MasterdataCategoryStatus("ACTIVE"), MasterdataCategoryStatusActive)
	assert.Equal(t, MasterdataCategoryStatus("INACTIVE"), MasterdataCategoryStatusInactive)
}

func TestMasterdataCategory_JSONTags(t *testing.T) {
	now := time.Now()
	description := "Test description"
	parentID := uuid.New()

	cat := MasterdataCategory{
		ID:                 uuid.New(),
		Code:               "TEST",
		Name:               "Test Category",
		Description:        &description,
		ParentCategoryID:   &parentID,
		IsSystem:           true,
		IsTenantExtensible: false,
		SortOrder:          1,
		Status:             MasterdataCategoryStatusActive,
		Version:            1,
		Timestamps: Timestamps{
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
