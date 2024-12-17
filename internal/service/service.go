package service

import (
	"awesomeProject/internal/store"
	"context"
)

type Service struct {
	Auth
}

func NewService(store *store.Store, jwtSecret, refreshSecret string) *Service {
	return &Service{
		Auth: NewAuthService(store.Auth, jwtSecret, refreshSecret),
	}
}

type Auth interface {
	GenerateTokens(ctx context.Context, userId string, clientIp string) (string, string, error)
	ExtractUserIDFromAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, string, error)
	Refresh(ctx context.Context, userID string, providedRefreshToken string, newClientIP string) (newAccess string, err error)
}
