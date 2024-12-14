package config

import (
	"awesomeProject/internal/logger"
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	JwtKet       string `env:"JWT_KET"`
	DbConnString string `env:"DB_CONNECTION_STRING"`
}

func InitConfig(ctx context.Context) *Config {
	log := logger.LoggerFromContext(ctx)
	var cfg Config

	err := cleanenv.ReadConfig("config/.env", &cfg)

	if err != nil {
		log.Errorw("Error reading config", zap.Error(err))
		return nil
	}

	return &cfg
}