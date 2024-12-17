package handler

import (
	"net/http"

	"github.com/docker/docker/daemon/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetTokens(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userId is empty",
		})
		return
	}

	clientIp := c.ClientIP()
	access, refresh, err := h.service.GenerateTokens(c.Request.Context(), userId, clientIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	log := logger.LoggerFromContext(c.Request.Context())
	var req struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("error with binding json", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh token is empty",
		})
		return
	}

	if req.AccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "access token is empty",
		})
		return
	}

	claims, userID, err := h.service.ExtractUserIDFromAccessToken(c.Request.Context(), req.AccessToken)
	if err != nil {
		log.Errorw("error with getting jwt claims", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	oldIPVal, ok := claims["client_ip"].(string)
	if !ok {
		log.Errorw("error with getting clientIp in claims", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	oldIP := oldIPVal

	newIP := c.Request.RemoteAddr

	newAccess, err := h.service.Refresh(c.Request.Context(), userID, req.RefreshToken, newIP)
	if err != nil {
		log.Errorw("error with refreshing tokens", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

}
