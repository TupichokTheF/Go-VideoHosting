package video

import (
	"errors"
	"fmt"
)

var (
	ErrInvalaidTitle      = errors.New("Invalid title of video")
	ErrInvalidDescription = errors.New("Invalid description of video")
)

type ValidationError struct {
	Field  string
	Reason string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Reason)
}
