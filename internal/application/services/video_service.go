package services

import "project/internal/domain/video"


type VideoService struct {
	videoRepo video.Repository
}

func NewVideoService(videoRepo video.Repository) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
	}
}

