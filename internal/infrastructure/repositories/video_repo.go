package repositories

import (
	"context"
	"fmt"
	"project/internal/domain/video"

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

func (repo *VideoRepository) AddVideo(ctx context.Context, video *video.Video) (int, error) {
	var id int
	state := video.State()
	err := repo.pool.QueryRow(ctx,
		`INSERT INTO videos(owner_id, title, description, status_id, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING video_id`, 
		state.OwnerID, state.Title, state.Description, state.FromStatusToID(), state.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error while adding video: %w", err)
	}

	return id, nil
}
