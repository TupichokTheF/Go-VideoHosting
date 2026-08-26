package app_errors

import "errors"

var (
	ErrInvalidToken = errors.New("Invalid refresh token")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrForbidden    = errors.New("Forbidden")
)
