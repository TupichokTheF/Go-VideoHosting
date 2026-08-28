package app_ports

import "time"

type Claims struct {
	UserID int
	JTI    string
	Exp    time.Time
}

type JWTManager interface {
	NewAccessToken(userID int) (string, error)
	NewRefreshToken(userID int) (string, error)
	ParseAccessToken(inputToken string) (*Claims, error)
	ParseRefreshToken(inputToken string) (*Claims, error)
}
