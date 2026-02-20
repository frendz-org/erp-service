package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID      uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions,omitempty"`
	BranchID    *uuid.UUID `json:"branch_id,omitempty"`
	SessionID   uuid.UUID  `json:"session_id"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return c.ExpiresAt.Before(time.Now())
}

func (c *JWTClaims) IsPlatformAdmin() bool {
	if c.TenantID != nil {
		return false
	}
	for _, role := range c.Roles {
		if role == "PLATFORM_ADMIN" {
			return true
		}
	}
	return false
}

func (c *JWTClaims) HasRole(roleCode string) bool {
	for _, role := range c.Roles {
		if role == roleCode {
			return true
		}
	}
	return false
}

func (c *JWTClaims) IsTenantUser() bool {
	return c.TenantID != nil
}

func (c *JWTClaims) GetTenantID() uuid.UUID {
	if c.TenantID == nil {
		return uuid.Nil
	}
	return *c.TenantID
}

func (c *JWTClaims) GetBranchID() uuid.UUID {
	if c.BranchID == nil {
		return uuid.Nil
	}
	return *c.BranchID
}

func (c *JWTClaims) GetProductID() uuid.UUID {
	if c.ProductID == nil {
		return uuid.Nil
	}
	return *c.ProductID
}

func (c *JWTClaims) HasProductContext() bool {
	return c.ProductID != nil
}

func (c *JWTClaims) HasPermission(permissionCode string) bool {
	for _, perm := range c.Permissions {
		if perm == permissionCode {
			return true
		}
	}
	return false
}

func (c *JWTClaims) HasAudience(audience string) bool {
	for _, aud := range c.Audience {
		if aud == audience {
			return true
		}
	}
	return false
}

type ProductClaim struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductCode string    `json:"product_code"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

type TenantClaim struct {
	TenantID uuid.UUID      `json:"tenant_id"`
	Products []ProductClaim `json:"products,omitempty"`
}

type MultiTenantClaims struct {
	UserID    uuid.UUID     `json:"user_id"`
	Email     string        `json:"email"`
	Roles     []string      `json:"roles,omitempty"`
	Tenants   []TenantClaim `json:"tenants,omitempty"`
	SessionID uuid.UUID     `json:"session_id"`
	jwt.RegisteredClaims
}

func (c *MultiTenantClaims) IsPlatformAdmin() bool {
	for _, role := range c.Roles {
		if role == "PLATFORM_ADMIN" {
			return true
		}
	}
	return false
}

func (c *MultiTenantClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return c.ExpiresAt.Before(time.Now())
}

func (c *MultiTenantClaims) HasTenant(tenantID uuid.UUID) bool {
	for _, t := range c.Tenants {
		if t.TenantID == tenantID {
			return true
		}
	}
	return false
}

func (c *MultiTenantClaims) GetTenantClaim(tenantID uuid.UUID) *TenantClaim {
	for i := range c.Tenants {
		if c.Tenants[i].TenantID == tenantID {
			return &c.Tenants[i]
		}
	}
	return nil
}

func (c *MultiTenantClaims) HasRoleInProduct(tenantID, productID uuid.UUID, roleCode string) bool {
	tc := c.GetTenantClaim(tenantID)
	if tc == nil {
		return false
	}
	for _, p := range tc.Products {
		if p.ProductID == productID {
			for _, r := range p.Roles {
				if r == roleCode {
					return true
				}
			}
		}
	}
	return false
}

func (c *MultiTenantClaims) HasPermissionInProduct(tenantID, productID uuid.UUID, permCode string) bool {
	tc := c.GetTenantClaim(tenantID)
	if tc == nil {
		return false
	}
	for _, p := range tc.Products {
		if p.ProductID == productID {
			for _, perm := range p.Permissions {
				if perm == permCode {
					return true
				}
			}
		}
	}
	return false
}
