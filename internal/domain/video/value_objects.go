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

func (title Title) String() string {
	return title.value
}

type Description struct{ value string }

func NewDescription(value string) (Description, error) {
	if utf8.RuneCountInString(value) > 1000 {
		return Description{}, ErrInvalidDescription
	}

	return Description{value: value}, nil
}

func (description Description) String() string {
	return description.value
}

type Status string

var (
	Draft    Status = "draft"
	Uploaded Status = "uploaded"
	Ready    Status = "ready"
	Deleted  Status = "deleted"
)
