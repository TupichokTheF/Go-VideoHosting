package video

import "time"

type Video struct {
	id          int
	title       Title
	description Description
	createdAt   time.Time
}

func New(inputTitle, inputDescription string, inputCreatedAt time.Time) (*Video, error) {
	title, err := NewTitle(inputTitle)
	if err != nil {
		return nil, err
	}

	description, err := NewDescription(inputDescription)
	if err != nil {
		return nil, err
	}

	return &Video{
		title: title,
		description: description,
		createdAt: inputCreatedAt,
	}, nil
}
