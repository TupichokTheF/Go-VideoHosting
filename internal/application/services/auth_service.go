package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	app_errors "project/internal/application/errors"
	app_ports "project/internal/application/ports"
	"project/internal/domain/user"
	"time"
)

type AuthService struct {
	userRepo   user.Repository
	jwtManager app_ports.JWTManager
	hasher     user.Hasher
	tokenCache app_ports.TokenCache
}

func NewAuthService(userRepo user.Repository,
	jwtManager app_ports.JWTManager,
	hasher user.Hasher,
	tokenCache app_ports.TokenCache) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		hasher:     hasher,
		tokenCache: tokenCache,
	}
}

func (authService *AuthService) RegisterUser(ctx context.Context, userCreateDTO *dtos.UserCreate) (*dtos.UserCreated, error) {
	createdUser, err := user.New(userCreateDTO.UserName, userCreateDTO.UserEmail, userCreateDTO.UserPassword, authService.hasher)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	userID, err := authService.userRepo.AddUser(ctx, createdUser)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	return &dtos.UserCreated{
		UserId: userID,
	}, nil
}

func (authService *AuthService) AuthorizeUser(ctx context.Context, authorizeDTO *dtos.Authorize) (*dtos.Tokens, error) {
	u, err := authService.userRepo.GetUserByUsername(ctx, authorizeDTO.Username)
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", err)
	}

	if ok := u.VerifyPassword(authorizeDTO.Password, authService.hasher); !ok {
		return nil, fmt.Errorf("user Authorization: %w", user.InvalidPassword)
	}

	accessToken, err := authService.jwtManager.NewAccessToken(u.ID())
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", app_errors.ErrInvalidToken)
	}

	refreshToken, err := authService.jwtManager.NewRefreshToken(u.ID())
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", app_errors.ErrInvalidToken)
	}

	return &dtos.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (authService *AuthService) RefreshToken(ctx context.Context, token string) (*dtos.Tokens, error) {
	tokenData, err := authService.jwtManager.ParseRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", app_errors.ErrInvalidToken)
	}

	ok, err := authService.tokenCache.IsRevoked(ctx, tokenData.JTI)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	
	if ok {
		return nil, app_errors.ErrTokenRevoked
	}

	accessToken, err := authService.jwtManager.NewAccessToken(tokenData.UserID)
	if err != nil {
		return nil, app_errors.ErrInvalidToken
	}

	return &dtos.Tokens{
		AccessToken: accessToken,
	}, nil
}

func (authService *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenData, err := authService.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	ttl := time.Until(tokenData.Exp)

	if err := authService.tokenCache.MarkAsRevoked(ctx, tokenData.JTI, ttl); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	return nil
}

func (authService *AuthService) Authenticate(ctx context.Context, accessToken string) (int, bool) {
	tokenData, err := authService.jwtManager.ParseAccessToken(accessToken)

	if err != nil {
		return 0, false
	}

	return tokenData.UserID, true
}
