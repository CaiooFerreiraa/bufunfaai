package errors

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyInUse   = errors.New("email already in use")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
)
