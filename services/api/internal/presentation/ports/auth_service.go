package pres_ports

import (
	"context"
	"project/internal/application/dtos"
)

type AuthService interface {
	RegisterUser(ctx context.Context, userCreateDTO *dtos.UserCreate) (*dtos.UserCreated, error)
	AuthorizeUser(ctx context.Context, authorizeDTO *dtos.Authorize) (*dtos.Tokens, error)
	RefreshToken(ctx context.Context, token string) (*dtos.Tokens, error)
	Logout(ctx context.Context, token string) error
	Authenticate(ctx context.Context, accessToken string) (int, bool)
}
