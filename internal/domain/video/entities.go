package video

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	id          uuid.UUID
	ownerID     int
	title       Title
	description Description
	status      Status
	size        int64
	createdAt   time.Time
}

func New(ownerID int, inputTitle, inputDescription string) (*Video, error) {
	title, err := NewTitle(inputTitle)
	if err != nil {
		return nil, &ValidationError{Field: "title", Reason: err.Error()}
	}

	description, err := NewDescription(inputDescription)
	if err != nil {
		return nil, &ValidationError{Field: "description", Reason: err.Error()}
	}

	return &Video{
		id:          uuid.Must(uuid.NewV7()),
		ownerID:     ownerID,
		title:       title,
		description: description,
		createdAt:   time.Now(),
		status:      Draft,
	}, nil
}

func (video *Video) ID() uuid.UUID {
	return video.id
}

func (video *Video) OwnerID() int {
	return video.ownerID
}

func (video *Video) Title() Title {
	return video.title
}

func (video *Video) Description() Description {
	return video.description
}

func (video *Video) Status() Status {
	return video.status
}

func (video *Video) CreatedAt() time.Time {
	return video.createdAt
}

func (video *Video) Size() int64 {
	return video.size
}

func (video *Video) MarkUploaded(inputSize int64) error {
	video.status = Uploaded
	video.size = inputSize

	return nil
}
