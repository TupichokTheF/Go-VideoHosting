package infra_ports


type JWTManagerInterface interface {
	NewAccessToken(userID int) (string, error)
	NewRefreshToken(userID int) (string, error)
	ParseAccessToken(inputToken string) (int, error)
	ParseRefreshToken(inputToken string) (int, error)
}

type HasherInterface interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}