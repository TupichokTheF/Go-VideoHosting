package infra_ports

import "context"


type TokenCacheInterface interface {
	SetRefreshToken(ctx context.Context, refresh string, userID int) error
	DeleteToken(ctx context.Context, userID int) error
	GetRefreshToken(ctx context.Context, userID int) (string, error) 
}