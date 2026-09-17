package user

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("User was not found")
	ErrAlreadyExist    = errors.New("User already exist")
	ErrInvalidPassword = errors.New("Invalid password")
)

type ValidationError struct {
	Field  string
	Reason string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Reason)
}
