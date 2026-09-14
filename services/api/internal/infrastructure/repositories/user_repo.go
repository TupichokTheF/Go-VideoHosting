package repositories

import (
	"context"
	"fmt"
	"project/internal/domain/user"

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
		return nil, fmt.Errorf("Get user %q: %w", username, user.NotFoundError)
	}

	return user.Reconstitute(u), nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, userID int) (*user.User, error) {
	var u *user.UserState = new(user.UserState)
	err := repo.pool.QueryRow(ctx,
		"SELECT user_id, username, email, password FROM users WHERE user_id = $1",
		userID).Scan(&u.ID, &u.UserName, &u.UserEmail, &u.UserPassword)

	if err != nil {
		return nil, fmt.Errorf("Get user %v: %w", userID, user.NotFoundError)
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
		fmt.Println(err)
		return 0, fmt.Errorf("add user: %w", user.AlreadyExistError)
	}

	return id, nil
}
