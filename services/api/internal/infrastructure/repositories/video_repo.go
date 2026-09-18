package repositories

import (
	"context"
	"errors"
	"fmt"
	"project/internal/domain/video"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (repo *VideoRepository) AddVideo(ctx context.Context, video *video.Video) (string, error) {
	var video_id string

	state := video.State()
	err := repo.pool.QueryRow(ctx,
		`INSERT INTO videos(video_id, owner_id, title, description, status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING video_id`,
		state.ID, state.OwnerID, state.Title, state.Description, state.Status, state.CreatedAt).Scan(&video_id)
	if err != nil {
		return "", fmt.Errorf("error while adding video %v : %w", video_id, err)
	}

	return video_id, nil
}

func (repo *VideoRepository) GetVideoByID(ctx context.Context, videoID uuid.UUID) (*video.Video, error) {
	var state video.State
	err := repo.pool.QueryRow(ctx,
		`SELECT video_id, owner_id, title, description, status, created_at 
		FROM videos WHERE video_id = $1`, videoID).
		Scan(&state.ID, &state.OwnerID, &state.Title, &state.Description, &state.Status, &state.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get video by id: %w", video.ErrNotFound)
		}
		return nil, fmt.Errorf("get video by id %v: %w", videoID, err)
	}

	return video.Reconstitute(&state), nil
}

func (repo *VideoRepository) UpdateVideo(ctx context.Context, v *video.Video) error {
	s := v.State()
	rows, err := repo.pool.Exec(ctx,
		`UPDATE videos SET status = $2, size = $3
		 WHERE video_id = $1`, s.ID, s.Status, s.Size)
	if err != nil {
		return fmt.Errorf("update video %s: %w", s.ID, err)
	}
	if rows.RowsAffected() == 0 {
		return fmt.Errorf("update video %s: %w", s.ID, video.ErrNotFound)
	}

	return nil
}
