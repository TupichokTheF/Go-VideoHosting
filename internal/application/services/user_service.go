package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	"project/internal/domain/user"
	infra_ports "project/internal/infrastructure/ports"
)

type UserService struct {
	userRepo   user.Repository
	tokenCache infra_ports.TokenCacheInterface
}

func NewUserService(userRepo user.Repository, tokenCache infra_ports.TokenCacheInterface) *UserService {
	return &UserService{
		userRepo:   userRepo,
		tokenCache: tokenCache,
	}
}

func (userService *UserService) GetUserInfo(ctx context.Context, userID int) (*dtos.UserInfoDTO, error) {
	u, err := userService.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &dtos.UserInfoDTO{
		UserID:    u.ID(),
		Username:  u.Username().String(),
		UserEmail: u.Email().String(),
	}, nil
}
