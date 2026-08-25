package video

import "time"

type State struct {
	ID          int
	OwnerID     int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
}

func (video *Video) State() *State {
	return &State{
		ID:          video.ID(),
		OwnerID:     video.OwnerID(),
		Title:       video.Title().String(),
		Description: video.Description().String(),
		Status:      string(video.Status()),
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
		createdAt:   state.CreatedAt,
	}
}

func (state *State) FromStatusToID() int {
	switch state.Status {
	case "draft":
		return 1
	case "uploaded":
		return 2
	case "ready":
		return 3
	case "deleted":
		return 4
	default:
		return -1
	}
}
