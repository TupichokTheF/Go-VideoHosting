package infra_ports

type JWTManager interface {
	NewAccessToken(userID int) (string, error)
	NewRefreshToken(userID int) (string, error)
	ParseAccessToken(inputToken string) (int, error)
	ParseRefreshToken(inputToken string) (int, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
