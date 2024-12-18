package handler

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/logger"
	"awesomeProject/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	jwtSecret string
	service   *service.Service
}

func NewHandler(service *service.Service, jwtSecret string) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret}
}

func (h *Handler) InitRoutes(ctx context.Context) http.Handler {
	log := logger.LoggerFromContext(ctx)
	router := gin.Default()
	router.Use(gin.Recovery(), logger.LoggerMiddleware(log))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("token", h.GetTokens)
			auth.Use(h.CheckAuth).POST("refresh", h.RefreshTokens)
		}
	}

	return router
}
