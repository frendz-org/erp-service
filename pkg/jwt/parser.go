package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseAccessToken(tokenString string, config *TokenConfig) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if config.SigningMethod == "RS256" {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected RS256)", token.Header["alg"])
			}
			return config.PublicKey, nil
		} else {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected HS256)", token.Header["alg"])
			}
			return []byte(config.AccessSecret), nil
		}
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrTokenSignature
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenInvalid
		}
		return nil, ErrTokenUnexpected
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.Issuer != config.Issuer {
			return nil, ErrTokenInvalid
		}

		if claims.IsExpired() {
			return nil, ErrTokenExpired
		}

		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func ParseMultiTenantAccessToken(tokenString string, config *TokenConfig) (*MultiTenantClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MultiTenantClaims{}, func(token *jwt.Token) (interface{}, error) {
		if config.SigningMethod == "RS256" {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected RS256)", token.Header["alg"])
			}
			return config.PublicKey, nil
		} else {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected HS256)", token.Header["alg"])
			}
			return []byte(config.AccessSecret), nil
		}
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrTokenSignature
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenInvalid
		}
		return nil, ErrTokenUnexpected
	}

	if claims, ok := token.Claims.(*MultiTenantClaims); ok && token.Valid {
		if claims.Issuer != config.Issuer {
			return nil, ErrTokenInvalid
		}

		if claims.IsExpired() {
			return nil, ErrTokenExpired
		}

		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func ParseRefreshToken(tokenString string, config *TokenConfig) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if config.SigningMethod == "RS256" {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected RS256)", token.Header["alg"])
			}
			return config.PublicKey, nil
		} else {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v (expected HS256)", token.Header["alg"])
			}
			return []byte(config.RefreshSecret), nil
		}
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrTokenSignature
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenInvalid
		}
		return nil, ErrTokenUnexpected
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Issuer != config.Issuer {
			return nil, ErrTokenInvalid
		}

		return claims, nil
	}

	return nil, ErrTokenInvalid
}
