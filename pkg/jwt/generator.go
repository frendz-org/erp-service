package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenConfig struct {
	SigningMethod string

	AccessSecret  string
	RefreshSecret string

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey

	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
	Audience      []string
}

func GenerateAccessToken(
	userID uuid.UUID,
	email string,
	tenantID *uuid.UUID,
	productID *uuid.UUID,
	roles []string,
	permissions []string,
	branchID *uuid.UUID,
	sessionID uuid.UUID,
	config *TokenConfig,
) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AccessExpiry)

	claims := &JWTClaims{
		UserID:      userID,
		Email:       email,
		TenantID:    tenantID,
		ProductID:   productID,
		Roles:       roles,
		Permissions: permissions,
		BranchID:    branchID,
		SessionID:   sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID.String(),
			Issuer:    config.Issuer,
			Audience:  config.Audience,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	var token *jwt.Token
	var signingKey interface{}

	if config.SigningMethod == "RS256" {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signingKey = config.PrivateKey
	} else {

		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signingKey = []byte(config.AccessSecret)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func GenerateMultiTenantAccessToken(
	userID uuid.UUID,
	email string,
	roles []string,
	tenants []TenantClaim,
	sessionID uuid.UUID,
	config *TokenConfig,
) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AccessExpiry)

	claims := &MultiTenantClaims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		Tenants:   tenants,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID.String(),
			Issuer:    config.Issuer,
			Audience:  config.Audience,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	var token *jwt.Token
	var signingKey interface{}

	if config.SigningMethod == "RS256" {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signingKey = config.PrivateKey
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signingKey = []byte(config.AccessSecret)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign multi-tenant token: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken(
	userID uuid.UUID,
	sessionID uuid.UUID,
	config *TokenConfig,
) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.RefreshExpiry)

	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    config.Issuer,
		Audience:  config.Audience,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(now),
		ID:        sessionID.String(),
	}

	var token *jwt.Token
	var signingKey interface{}

	if config.SigningMethod == "RS256" {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signingKey = config.PrivateKey
	} else {

		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signingKey = []byte(config.RefreshSecret)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}
