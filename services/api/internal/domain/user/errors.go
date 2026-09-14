package user

import (
	"errors"
	"fmt"
)

var (
	NotFoundError     = errors.New("User was not found")
	AlreadyExistError = errors.New("User already exist")
	InvalidPassword   = errors.New("Invalid password")
)

type ValidationError struct {
	Field  string
	Reason string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Reason)
}
