package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SAMLAttributeMapping struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Roles string `json:"roles"`
}

type SAMLConfiguration struct {
	ID                 uuid.UUID       `json:"id" db:"id"`
	TenantID           uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	IDPEntityID        string          `json:"idp_entity_id" db:"idp_entity_id"`
	IDPSSOURL          string          `json:"idp_sso_url" db:"idp_sso_url"`
	IDPSLOURL          *string         `json:"idp_slo_url,omitempty" db:"idp_slo_url"`
	IDPCertificate     string          `json:"idp_certificate" db:"idp_certificate"`
	SPEntityID         string          `json:"sp_entity_id" db:"sp_entity_id"`
	SPACSURL           string          `json:"sp_acs_url" db:"sp_acs_url"`
	SPSLOURL           *string         `json:"sp_slo_url,omitempty" db:"sp_slo_url"`
	AttributeMapping   json.RawMessage `json:"attribute_mapping" db:"attribute_mapping"`
	RoleMapping        json.RawMessage `json:"role_mapping,omitempty" db:"role_mapping"`
	AutoProvisionUsers bool            `json:"auto_provision_users" db:"auto_provision_users"`
	DefaultBranchID    *uuid.UUID      `json:"default_branch_id,omitempty" db:"default_branch_id"`
	IsActive           bool            `json:"is_active" db:"is_active"`
	Timestamps
}

func (s *SAMLConfiguration) GetAttributeMapping() (*SAMLAttributeMapping, error) {
	var mapping SAMLAttributeMapping
	if err := json.Unmarshal(s.AttributeMapping, &mapping); err != nil {
		return nil, err
	}
	return &mapping, nil
}

func (s *SAMLConfiguration) GetRoleMapping() (map[string]string, error) {
	var mapping map[string]string
	if err := json.Unmarshal(s.RoleMapping, &mapping); err != nil {
		return nil, err
	}
	return mapping, nil
}

func (s *SAMLConfiguration) SetAttributeMapping(mapping *SAMLAttributeMapping) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		return err
	}
	s.AttributeMapping = data
	return nil
}

func (s *SAMLConfiguration) SetRoleMapping(mapping map[string]string) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		return err
	}
	s.RoleMapping = data
	return nil
}

func DefaultSAMLAttributeMapping() *SAMLAttributeMapping {
	return &SAMLAttributeMapping{
		Email: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
		Name:  "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name",
		Roles: "http://schemas.microsoft.com/ws/2008/06/identity/claims/role",
	}
}

func NewSAMLConfiguration(tenantID uuid.UUID) *SAMLConfiguration {
	defaultMapping := DefaultSAMLAttributeMapping()
	mappingJSON, _ := json.Marshal(defaultMapping)
	roleMapping, _ := json.Marshal(map[string]string{})

	now := time.Now()
	return &SAMLConfiguration{
		ID:                 uuid.New(),
		TenantID:           tenantID,
		AttributeMapping:   mappingJSON,
		RoleMapping:        roleMapping,
		AutoProvisionUsers: true,
		IsActive:           true,
		Timestamps: Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
