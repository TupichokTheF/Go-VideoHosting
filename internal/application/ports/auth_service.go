package app_ports

import (
	"context"
	"project/internal/application/dtos"
)

type AuthService interface {
	RegisterUser(ctx context.Context, userCreateDTO *dtos.UserCreateDTO) (*dtos.UserCreatedDTO, error)
	AuthorizeUser(ctx context.Context, authorizeDTO *dtos.AuthorizeDTO) (*dtos.TokensDTO, error)
	RefreshToken(ctx context.Context, token string) (*dtos.TokensDTO, error)
	Logout(ctx context.Context, token string) error
	IsAuthorized(ctx context.Context, accessToken string) (int, bool)
	IsLoggedOut(ctx context.Context, refreshToken string) bool
}
