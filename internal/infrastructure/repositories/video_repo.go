package repositories

import (
	"context"
	"fmt"
	"project/internal/domain/video"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoRepository struct {
	pool *pgxpool.Pool
}

func NewVideoRepository(pool *pgxpool.Pool) *VideoRepository {
	return &VideoRepository{
		pool: pool,
	}
}

func (repo *VideoRepository) AddVideo(ctx context.Context, video *video.Video) error {
	var id int
	state := video.State()
	err := repo.pool.QueryRow(ctx,
		`INSERT INTO videos(video_id, owner_id, title, description, status, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING video_id`,
		state.ID, state.OwnerID, state.Title, state.Description, state, state.CreatedAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("error while adding video: %w", err)
	}

	return nil
}

func (repo *VideoRepository) GetVideoByID(ctx context.Context, videoID uuid.UUID) (*video.Video, error) {
	var state video.State
	err := repo.pool.QueryRow(ctx,
		`SELECT video_id, owner_id, title, description, status, created_at 
		FROM videos WHERE video_id = $1`, videoID).
		Scan(&state.ID, &state.OwnerID, &state.Title, &state.Description, &state.Status, &state.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get video by id: %w", video.ErrNotFound)
	}

	return video.Reconstitute(&state), nil
}

func (repo *VideoRepository) UpdateVideo(ctx context.Context, video *video.Video) error {
	s := video.State()
	_, err := repo.pool.Exec(ctx,
		`UPDATE videos SET status = $2, size_bytes = $3
		 WHERE video_id = $1`, s.ID, s.Status, s.Size)
	if err != nil {
		return fmt.Errorf("update video %s: %w", s.ID, err)
	}

	return nil
}
