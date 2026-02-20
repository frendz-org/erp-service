package iamclient

import (
	"fmt"

	"iam-service/pkg/jwt"
)

type Client struct {
	AccessSecret     string
	Issuer           string
	RequiredAudience string
	RequiredProduct  string
}

type Config struct {
	AccessSecret     string
	Issuer           string
	RequiredAudience string
	RequiredProduct  string
}

func NewClient(config *Config) (*Client, error) {
	if config.AccessSecret == "" {
		return nil, fmt.Errorf("AccessSecret is required")
	}

	return &Client{
		AccessSecret:     config.AccessSecret,
		Issuer:           config.Issuer,
		RequiredAudience: config.RequiredAudience,
		RequiredProduct:  config.RequiredProduct,
	}, nil
}

func (c *Client) ValidateToken(tokenString string) (*jwt.JWTClaims, error) {
	tokenConfig := &jwt.TokenConfig{
		SigningMethod: "HS256",
		AccessSecret: c.AccessSecret,
		Issuer:       c.Issuer,
	}

	claims, err := jwt.ParseAccessToken(tokenString, tokenConfig)
	if err != nil {
		return nil, err
	}

	if c.RequiredAudience != "" && !claims.HasAudience(c.RequiredAudience) {
		return nil, fmt.Errorf("invalid audience: token not intended for this service")
	}

	if c.RequiredProduct != "" {
		if claims.ProductID == nil {
			return nil, fmt.Errorf("product context required but not found in token")
		}
	}

	return claims, nil
}

func (c *Client) HasPermission(claims *jwt.JWTClaims, permissionCode string) bool {
	return claims.HasPermission(permissionCode)
}

func (c *Client) HasRole(claims *jwt.JWTClaims, roleCode string) bool {
	return claims.HasRole(roleCode)
}

func (c *Client) HasAnyRole(claims *jwt.JWTClaims, roleCodes []string) bool {
	for _, role := range roleCodes {
		if claims.HasRole(role) {
			return true
		}
	}
	return false
}

func (c *Client) HasAllRoles(claims *jwt.JWTClaims, roleCodes []string) bool {
	for _, role := range roleCodes {
		if !claims.HasRole(role) {
			return false
		}
	}
	return true
}
