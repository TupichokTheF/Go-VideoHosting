package video

import (
	"strings"
	"unicode/utf8"
)

type Title struct{ value string }

func NewTitle(value string) (Title, error) {
	if strings.TrimSpace(value) == "" {
		return Title{}, ErrInvalaidTitle
	}

	return Title{value: value}, nil
}

type Description struct{ value string }

func NewDescription(value string) (Description, error) {
	if utf8.RuneCountInString(value) > 1000 {
		return Description{}, ErrInvalidDescription
	}

	return Description{value: value}, nil
}
