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

	serviceLevel := service.NewService(storeLevel, cfg.JwtKet)
	handlerLevel := handler.NewHandler(serviceLevel)

	//todo graceful shutdown
	if err := server.InitServer("8081", handlerLevel.InitRoutes(ctx)); err != nil {
		log.Errorw("error with initing server", zap.Error(err))
		return
	}
}
