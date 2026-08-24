package video

import "errors"

var (
	ErrInvalaidTitle      = errors.New("Invalid title of video")
	ErrInvalidDescription = errors.New("Invalid description of video")
)
