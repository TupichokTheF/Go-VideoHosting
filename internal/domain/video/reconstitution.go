package video

import (
	"time"

	"github.com/google/uuid"
)

type State struct {
	ID          uuid.UUID
	OwnerID     int
	Title       string
	Description string
	Status      string
	Size        int64
	CreatedAt   time.Time
}

func (video *Video) State() *State {
	return &State{
		ID:          video.ID(),
		OwnerID:     video.OwnerID(),
		Title:       video.Title().String(),
		Description: video.Description().String(),
		Status:      string(video.Status()),
		Size:        video.Size(),
		CreatedAt:   video.CreatedAt(),
	}
}

func Reconstitute(state *State) *Video {
	return &Video{
		id:          state.ID,
		ownerID:     state.OwnerID,
		title:       Title{value: state.Title},
		description: Description{value: state.Description},
		status:      Status(state.Status),
		size:        state.Size,
		createdAt:   state.CreatedAt,
	}
}
