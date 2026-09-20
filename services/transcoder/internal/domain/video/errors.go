package video

import (
	"errors"
	"fmt"
)

var (
	ErrInvalaidTitle      = errors.New("Invalid title of video")
	ErrInvalidDescription = errors.New("Invalid description of video")
	ErrNotFound           = errors.New("Video doesn't exist")
	ErrAlreadyExist       = errors.New("Video already exist")
	ErrInvalidTransition  = errors.New("Invalid transition of video")
	ErrNotAvailable       = errors.New("Video is not available to show")
)

type ValidationError struct {
	Field  string
	Reason string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Reason)
}
