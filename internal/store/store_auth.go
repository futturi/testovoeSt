package store

import (
	"awesomeProject/internal/entites"
	"awesomeProject/internal/logger"
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthStore struct {
	DB *gorm.DB
}

func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{
		DB: db,
	}
}

func (r *AuthStore) InsertUserInfo(ctx context.Context, userId string, clientIp string, token string) error {
	log := logger.LoggerFromContext(ctx)

	tx := r.DB.Create(&entites.RefreshToken{
		UserID:   userId,
		Hash:     token,
		IssuedAt: time.Now(),
		ClientIp: clientIp,
	})

	if tx.Error != nil {
		log.Errorw("error with inserting refresh token", zap.Error(tx.Error))
		tx.Rollback()
		return tx.Error
	}

	return nil
}

func (r *AuthStore) GetUserById(ctx context.Context, userId string) (*entites.User, error) {
	var user entites.User

	log := logger.LoggerFromContext(ctx)

	tx := r.DB.Where("id = ?", userId).First(&user)

	if tx.Error != nil {
		log.Errorw("error with getting user by id", zap.Error(tx.Error))
		return nil, tx.Error
	}

	return &user, nil
}
