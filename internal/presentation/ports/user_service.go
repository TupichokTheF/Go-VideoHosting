package app_ports

import (
	"context"
	"project/internal/application/dtos"
)

type UserService interface {
	GetUserInfo(ctx context.Context, userID int) (*dtos.UserInfo, error)
}
