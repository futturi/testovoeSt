package store

import (
	"awesomeProject/internal/entites"
	"awesomeProject/internal/logger"
	"context"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func InitDB(ctx context.Context, connString string) (*Store, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debugw("connecting to database")
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Errorw("error with opening database conn", zap.Error(err))
		return nil, err
	}

	if err := db.AutoMigrate(&entites.User{}); err != nil {
		log.Errorw("error with migrating database scheme", zap.Error(err))
		return nil, err
	}

	log.Debug("database is connected")

	return &Store{
		db: db,
	}, nil
}
