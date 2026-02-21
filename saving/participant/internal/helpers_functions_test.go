package internal

import (
	"testing"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateParticipantOwnership(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	otherTenantID := uuid.New()
	otherProductID := uuid.New()

	tests := []struct {
		name        string
		participant *entity.Participant
		tenantID    uuid.UUID
		productID   uuid.UUID
		wantErr     bool
	}{
		{
			name: "success - matching tenant and product",
			participant: &entity.Participant{
				TenantID:  tenantID,
				ProductID: productID,
			},
			tenantID:  tenantID,
			productID: productID,
			wantErr:   false,
		},
		{
			name: "error - different tenant (BOLA)",
			participant: &entity.Participant{
				TenantID:  tenantID,
				ProductID: productID,
			},
			tenantID:  otherTenantID,
			productID: productID,
			wantErr:   true,
		},
		{
			name: "error - different product (BOLA)",
			participant: &entity.Participant{
				TenantID:  tenantID,
				ProductID: productID,
			},
			tenantID:  tenantID,
			productID: otherProductID,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParticipantOwnership(tt.participant, tt.tenantID, tt.productID)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, errors.KindForbidden, appErr.Kind)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEditableState(t *testing.T) {
	tests := []struct {
		name        string
		participant *entity.Participant
		wantErr     bool
	}{
		{
			name: "success - DRAFT participant is editable",
			participant: &entity.Participant{
				Status: entity.ParticipantStatusDraft,
			},
			wantErr: false,
		},
		{
			name: "success - REJECTED participant is editable",
			participant: &entity.Participant{
				Status: entity.ParticipantStatusRejected,
			},
			wantErr: false,
		},
		{
			name: "error - PENDING_APPROVAL participant is not editable",
			participant: &entity.Participant{
				Status: entity.ParticipantStatusPendingApproval,
			},
			wantErr: true,
		},
		{
			name: "error - APPROVED participant is not editable",
			participant: &entity.Participant{
				Status: entity.ParticipantStatusApproved,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEditableState(tt.participant)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, errors.KindBadRequest, appErr.Kind)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		want      string
	}{
		{
			name:      "allowed field - ktp_photo",
			fieldName: "ktp_photo",
			want:      "ktp_photo",
		},
		{
			name:      "allowed field - passport_photo",
			fieldName: "passport_photo",
			want:      "passport_photo",
		},
		{
			name:      "path traversal attempt - ../ prefix",
			fieldName: "../etc/passwd",
			want:      "unknown",
		},
		{
			name:      "path traversal attempt - absolute path",
			fieldName: "/etc/passwd",
			want:      "unknown",
		},
		{
			name:      "disallowed field - arbitrary name",
			fieldName: "malicious_file",
			want:      "unknown",
		},
		{
			name:      "dot only",
			fieldName: ".",
			want:      "unknown",
		},
		{
			name:      "dot dot only",
			fieldName: "..",
			want:      "unknown",
		},
		{
			name:      "contains slash",
			fieldName: "ktp/photo",
			want:      "unknown",
		},
		{
			name:      "contains backslash",
			fieldName: "ktp\\photo",
			want:      "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFieldName(tt.fieldName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "normal filename",
			filename: "document.pdf",
			want:     "document.pdf",
		},
		{
			name:     "path traversal attempt - ../ prefix",
			filename: "../../../etc/passwd",
			want:     "passwd",
		},
		{
			name:     "path traversal attempt - absolute path",
			filename: "/var/log/secrets.txt",
			want:     "secrets.txt",
		},
		{
			name:     "dot only",
			filename: ".",
			want:     "upload",
		},
		{
			name:     "dot dot only",
			filename: "..",
			want:     "upload",
		},
		{
			name:     "contains slash in middle",
			filename: "path/to/file.jpg",
			want:     "file.jpg",
		},
		{
			name:     "contains backslash",
			filename: "path\\to\\file.jpg",
			want:     "upload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateObjectKey(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	tests := []struct {
		name          string
		tenantID      uuid.UUID
		productID     uuid.UUID
		participantID uuid.UUID
		fieldName     string
		filename      string
		wantPrefix    string
		wantContains  []string
	}{
		{
			name:          "normal upload",
			tenantID:      tenantID,
			productID:     productID,
			participantID: participantID,
			fieldName:     "ktp_photo",
			filename:      "ktp.jpg",
			wantPrefix:    "participants/",
			wantContains:  []string{tenantID.String(), productID.String(), participantID.String(), "ktp_photo", "ktp.jpg"},
		},
		{
			name:          "sanitizes malicious field name",
			tenantID:      tenantID,
			productID:     productID,
			participantID: participantID,
			fieldName:     "../etc/passwd",
			filename:      "file.jpg",
			wantPrefix:    "participants/",
			wantContains:  []string{tenantID.String(), productID.String(), participantID.String(), "unknown", "file.jpg"},
		},
		{
			name:          "sanitizes malicious filename",
			tenantID:      tenantID,
			productID:     productID,
			participantID: participantID,
			fieldName:     "ktp_photo",
			filename:      "/etc/passwd",
			wantPrefix:    "participants/",
			wantContains:  []string{tenantID.String(), productID.String(), participantID.String(), "ktp_photo", "passwd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateObjectKey(tt.tenantID, tt.productID, tt.participantID, tt.fieldName, tt.filename)

			assert.True(t, len(got) > 0, "object key should not be empty")
			if tt.wantPrefix != "" {
				assert.Contains(t, got, tt.wantPrefix)
			}
			for _, contains := range tt.wantContains {
				assert.Contains(t, got, contains)
			}
		})
	}
}
