package jwt

import "errors"

var (
	ErrTokenInvalid    = errors.New("token is invalid")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenMalformed  = errors.New("token is malformed")
	ErrTokenSignature  = errors.New("token signature is invalid")
	ErrTokenUnexpected = errors.New("unexpected error parsing token")
)
