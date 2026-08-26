package dtos

import "github.com/google/uuid"

type PresignedURL struct {
	VideoID string
	URL     string
}

type CreateVideo struct {
	OwnerID     int
	Title       string
	Description string
}

type GetVideo struct {
	VideoID string
}

type CompleteVideo struct {
	VideoID uuid.UUID
	UserID  int
}
