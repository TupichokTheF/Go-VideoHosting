package services

import (
	"context"
	"fmt"
	"project/internal/application/dtos"
	"project/internal/domain/user"
)

type UserService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (userService *UserService) GetUserInfo(ctx context.Context, userID int) (*dtos.UserInfo, error) {
	u, err := userService.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &dtos.UserInfo{
		UserID:    u.ID(),
		Username:  u.Username().String(),
		UserEmail: u.Email().String(),
	}, nil
}
