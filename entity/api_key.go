package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

type AdminAPIKey struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	TenantID      uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	KeyName       string          `json:"key_name" db:"key_name"`
	KeyHash       string          `json:"-" db:"key_hash"`
	KeyPrefix     string          `json:"key_prefix" db:"key_prefix"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty" db:"created_by"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	ExpiresAt     *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	RevokedAt     *time.Time      `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokedBy     *uuid.UUID      `json:"revoked_by,omitempty" db:"revoked_by"`
	RevokedReason *string         `json:"revoked_reason,omitempty" db:"revoked_reason"`
	LastUsedAt    *time.Time      `json:"last_used_at,omitempty" db:"last_used_at"`
	LastUsedIP    net.IP          `json:"last_used_ip,omitempty" db:"last_used_ip"`
	IPWhitelist   json.RawMessage `json:"ip_whitelist,omitempty" db:"ip_whitelist"`
	IsActive      bool            `json:"is_active" db:"is_active"`
}

func (a *AdminAPIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

func (a *AdminAPIKey) IsRevoked() bool {
	return a.RevokedAt != nil
}

func (a *AdminAPIKey) IsValid() bool {
	return a.IsActive && !a.IsExpired() && !a.IsRevoked()
}

func (a *AdminAPIKey) GetIPWhitelist() ([]string, error) {
	var whitelist []string
	if a.IPWhitelist == nil {
		return whitelist, nil
	}
	if err := json.Unmarshal(a.IPWhitelist, &whitelist); err != nil {
		return nil, err
	}
	return whitelist, nil
}

func (a *AdminAPIKey) SetIPWhitelist(whitelist []string) error {
	data, err := json.Marshal(whitelist)
	if err != nil {
		return err
	}
	a.IPWhitelist = data
	return nil
}

func (a *AdminAPIKey) IsIPAllowed(ip net.IP) bool {
	whitelist, err := a.GetIPWhitelist()
	if err != nil || len(whitelist) == 0 {
		return true
	}

	ipStr := ip.String()
	for _, allowedIP := range whitelist {

		if allowedIP == ipStr {
			return true
		}

		_, cidr, err := net.ParseCIDR(allowedIP)
		if err == nil && cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func (a *AdminAPIKey) Revoke(revokedBy uuid.UUID, reason string) {
	now := time.Now()
	a.RevokedAt = &now
	a.RevokedBy = &revokedBy
	a.RevokedReason = &reason
	a.IsActive = false
}

func (a *AdminAPIKey) UpdateLastUsed(ip net.IP) {
	now := time.Now()
	a.LastUsedAt = &now
	a.LastUsedIP = ip
}

func NewAdminAPIKey(tenantID uuid.UUID, keyName, keyHash, keyPrefix string, createdBy *uuid.UUID) *AdminAPIKey {
	emptyWhitelist, _ := json.Marshal([]string{})
	return &AdminAPIKey{
		ID:          uuid.New(),
		TenantID:    tenantID,
		KeyName:     keyName,
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		IPWhitelist: emptyWhitelist,
		IsActive:    true,
	}
}

func NewAdminAPIKeyWithExpiry(tenantID uuid.UUID, keyName, keyHash, keyPrefix string, createdBy *uuid.UUID, expiresAt time.Time) *AdminAPIKey {
	key := NewAdminAPIKey(tenantID, keyName, keyHash, keyPrefix, createdBy)
	key.ExpiresAt = &expiresAt
	return key
}
