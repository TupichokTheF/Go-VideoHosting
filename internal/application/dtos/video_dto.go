package dtos

type PresignedURL struct {
	URL string
}

type CreateVideo struct {
	OwnerID     int
	Title       string
	Description string
}

type GetVideo struct {
	VideoID int
}
