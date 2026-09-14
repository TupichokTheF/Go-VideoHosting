package security

import (
	"errors"
	"fmt"
	app_ports "project/internal/application/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(accessSecret []byte, refreshSecret []byte, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (manager *JWTManager) NewAccessToken(userID int) (string, error) {
	return manager.newToken(manager.accessSecret, userID, manager.accessTTL)
}

func (manager *JWTManager) NewRefreshToken(userID int) (string, error) {
	return manager.newToken(manager.refreshSecret, userID, manager.refreshTTL)
}

func (manager *JWTManager) newToken(secret []byte, userID int, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": uuid.NewString(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func (m *JWTManager) ParseRefreshToken(inputToken string) (*app_ports.Claims, error) {
	return m.parseToken(inputToken, m.refreshSecret)
}

func (m *JWTManager) ParseAccessToken(inputToken string) (*app_ports.Claims, error) {
	return m.parseToken(inputToken, m.accessSecret)
}

func (m *JWTManager) parseToken(inputToken string, secret []byte) (*app_ports.Claims, error) {
	token, err := jwt.Parse(inputToken, func(t *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return m.parseClaims(claims)
}

func (m *JWTManager) parseClaims(claims jwt.MapClaims) (*app_ports.Claims, error) {
	userID, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("invalid type of userID")
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("get exp: %w", err)
	}
	if expTime == nil {
		return nil, errors.New("missing exp claim")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid type of jti")
	}

	return &app_ports.Claims{
		UserID: int(userID),
		Exp:    expTime.Time,
		JTI:    jti,
	}, nil
}
