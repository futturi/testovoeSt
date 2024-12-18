package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/handler"
	"awesomeProject/internal/logger"
	"awesomeProject/internal/server"
	"awesomeProject/internal/service"
	"awesomeProject/internal/store"
	"context"

	"go.uber.org/zap"
)

// @title Auth API
// @description API 4 Auth Testing

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	log := logger.InitLogger()
	ctx = logger.ContextWithLogger(ctx, log)

	cfg := config.InitConfig(ctx)
	storeLevel, err := store.InitDB(ctx, cfg.DbConnString)
	if err != nil {
		log.Errorw("failed to connect to db", zap.Error(err))
		return
	}

	serviceLevel := service.NewService(storeLevel, cfg.JwtKet, cfg.RefreshKey)
	handlerLevel := handler.NewHandler(serviceLevel, cfg.JwtKet)

	//todo graceful shutdown
	if err := server.InitServer("8080", handlerLevel.InitRoutes(ctx)); err != nil {
		log.Errorw("error with initing server", zap.Error(err))
		return
	}
}
