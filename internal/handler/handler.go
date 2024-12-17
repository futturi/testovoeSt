package handler

import (
	"awesomeProject/internal/logger"
	"awesomeProject/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(ctx context.Context) http.Handler {
	log := logger.LoggerFromContext(ctx)
	router := gin.Default()
	router.Use(gin.Recovery(), logger.LoggerMiddleware(log))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("token", h.GetTokens)
		}
	}

	return router
}
