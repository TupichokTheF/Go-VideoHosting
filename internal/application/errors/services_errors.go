package app_errors

import "errors"

var (
	InvalidTokenError = errors.New("Invalid refresh token")
	UnauthorizedError = errors.New("Unauthorized")
)
