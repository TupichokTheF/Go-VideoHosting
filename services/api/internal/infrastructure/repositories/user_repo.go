package repositories

import (
	"context"
	"errors"
	"fmt"
	"project/internal/domain/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (repo *UserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var u *user.UserState = new(user.UserState)
	err := repo.pool.QueryRow(ctx,
		"SELECT user_id, username, email, password FROM users WHERE username = $1",
		username).Scan(&u.ID, &u.UserName, &u.UserEmail, &u.UserPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get user %q: %w", username, user.ErrNotFound)
		}
		return nil, fmt.Errorf("get user %q: %w", username, err)
	}

	return user.Reconstitute(u), nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, userID int) (*user.User, error) {
	var u *user.UserState = new(user.UserState)
	err := repo.pool.QueryRow(ctx,
		"SELECT user_id, username, email, password FROM users WHERE user_id = $1",
		userID).Scan(&u.ID, &u.UserName, &u.UserEmail, &u.UserPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get user %v: %w", userID, user.ErrNotFound)
		}
		return nil, fmt.Errorf("get user %v: %w", userID, err)
	}

	return user.Reconstitute(u), nil
}

func (repo *UserRepository) AddUser(ctx context.Context, inputUser *user.User) (int, error) {
	var id int
	userState := inputUser.State()
	err := repo.pool.QueryRow(ctx,
		`INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3)
		RETURNING user_id`,
		userState.UserName, userState.UserEmail, userState.UserPassword).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("add user: %w", user.ErrAlreadyExist)
		}
		return 0, fmt.Errorf("add user: %w", err)
	}

	return id, nil
}
