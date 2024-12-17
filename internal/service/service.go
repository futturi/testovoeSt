package service

import (
	"awesomeProject/internal/store"
	"context"
)

type Service struct {
	Auth
}

func NewService(store *store.Store, jwtSecret string) *Service {
	return &Service{
		Auth: NewAuthService(store.Auth, jwtSecret),
	}
}

type Auth interface {
	GenerateTokens(ctx context.Context, userId string, clientIp string) (string, string, error)
	ExtractUserIDFromAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, string, error)
	Refresh(ctx context.Context, userID string, providedRefreshToken string, newClientIP string) (newAccess string, err error)
}
