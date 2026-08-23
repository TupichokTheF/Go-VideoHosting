package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	app_errors "project/internal/application/errors"
	"project/internal/domain/user"
	infra_ports "project/internal/infrastructure/ports"
)


type AuthService struct {
	userRepo   user.Repository
	jwtManager infra_ports.JWTManagerInterface
	hasher     infra_ports.HasherInterface
	tokenCache infra_ports.TokenCacheInterface
}


func NewAuthService(userRepo user.Repository, 
	jwtManager infra_ports.JWTManagerInterface, 
	hasher infra_ports.HasherInterface,
	tokenCache infra_ports.TokenCacheInterface) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		hasher:     hasher,
		tokenCache: tokenCache,
	}
}

func (authService *AuthService) RegisterUser(ctx context.Context, userCreateDTO *dtos.UserCreateDTO) (*dtos.UserCreatedDTO, error) {
	createdUser, err := user.New(userCreateDTO.UserName, userCreateDTO.UserEmail, userCreateDTO.UserPassword, authService.hasher)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	userID, err := authService.userRepo.AddUser(ctx, createdUser)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	return &dtos.UserCreatedDTO{
		UserId: userID,
	}, nil
}

func (authService *AuthService) AuthorizeUser(ctx context.Context, authorizeDTO *dtos.AuthorizeDTO) (*dtos.TokensDTO, error) {
	u, err := authService.userRepo.GetUserByUsername(ctx, authorizeDTO.Username)
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", err)
	}

	if ok := u.VerifyPassword(authorizeDTO.Password, authService.hasher); !ok {
		return nil, fmt.Errorf("user Authorization: %w", user.InvalidPassword)
	}

	accessToken, err := authService.jwtManager.NewAccessToken(u.ID())
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", app_errors.InvalidTokenError)
	}

	refreshToken, err := authService.jwtManager.NewRefreshToken(u.ID())
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", app_errors.InvalidTokenError)
	}

	if err := authService.tokenCache.SetRefreshToken(ctx, refreshToken, u.ID()); err != nil {
		return nil, fmt.Errorf("user Authorization: %w", err)
	}

	return &dtos.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (authService *AuthService) RefreshToken(ctx context.Context, token string) (*dtos.TokensDTO, error) {
	userID, err := authService.jwtManager.ParseRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", app_errors.InvalidTokenError)
	}

	accessToken, err := authService.jwtManager.NewAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("user authorization: %w", app_errors.InvalidTokenError)
	}

	return &dtos.TokensDTO{
		AccessToken:  accessToken,
	}, nil
}

func (authService *AuthService) Logout(ctx context.Context, refreshToken string) error {
	userID, err := authService.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("logout: %w", app_errors.InvalidTokenError)
	}

	if err := authService.tokenCache.DeleteToken(ctx, userID); err != nil {
		return fmt.Errorf("logout: %w", app_errors.InvalidTokenError) 
	}

	return nil
}

func (authService *AuthService) IsAuthorized(ctx context.Context, accessToken string) (int, bool) {
	userID, err := authService.jwtManager.ParseAccessToken(accessToken)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}

	return userID, true
}

func (authService *AuthService) IsLoggedOut(ctx context.Context, refreshToken string) bool {
	userID, err := authService.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return true
	}

	token, err := authService.tokenCache.GetRefreshToken(ctx, userID)
	if err != nil || token == "" {
		return true
	}

	return false
}