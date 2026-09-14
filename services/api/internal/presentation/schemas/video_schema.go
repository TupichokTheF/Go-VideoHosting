package schemas

type CreateVideo struct {
	OwnerID     int    `json:"owner_id" example:"1"`
	Title       string `json:"title" example:"Cool video about cats"`
	Description string `json:"description" example:"Description of cool video about cat"`
}

type GetVideo struct {
	VideoID string `json:"video_id"`
}

type PresignedURL struct {
	VideoID string `json:"video_id"`
	URL     string `json:"url"`
}
