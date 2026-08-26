package infra_ports

type JWTManager interface {
	NewAccessToken(userID int) (string, error)
	NewRefreshToken(userID int) (string, error)
	ParseAccessToken(inputToken string) (int, error)
	ParseRefreshToken(inputToken string) (int, error)
}