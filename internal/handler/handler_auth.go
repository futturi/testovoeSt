package handler

import (
	"awesomeProject/internal/entites"
	"awesomeProject/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary GetTokens
// @Tags auth
// @Description gettokens
// @ID gettokens
// @Accept json
// @Produce json
// @Param userId query string true "userId"
// @Success 200 {object} entites.Response
// @Failure 400 {object} entites.Error
// @Failure 500 {object} entites.Error
// @Router /api/auth/token [post]
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

	c.JSON(http.StatusOK, entites.Response{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// @Summary RefreshToken
// @Tags auth
// @Security ApiKeyAuth
// @Description refreshTokens
// @ID refreshTokens
// @Accept json
// @Produce json
// @Param input body entites.RefreshRequest true "refreshToken"
// @Success 200 {object} entites.Response
// @Failure 400 {object} entites.Error
// @Failure 500 {object} entites.Error
// @Router /api/auth/refresh [post]
func (h *Handler) RefreshTokens(c *gin.Context) {
	log := logger.LoggerFromContext(c.Request.Context())
	var req entites.RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorw("error with binding json", zap.Error(err))
		c.JSON(http.StatusBadRequest, entites.Error{
			Error: err.Error(),
		})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, entites.Error{
			Error: "refresh token is empty",
		})
		return
	}

	accessToken := c.GetString("accessToken")

	if accessToken == "" {
		c.JSON(http.StatusBadRequest, entites.Error{
			Error: "access token is empty",
		})
		return
	}

	_, userID, err := h.service.ExtractUserIDFromAccessToken(c.Request.Context(), accessToken)
	if err != nil {
		log.Errorw("error with getting jwt claims", zap.Error(err))
		c.JSON(http.StatusBadRequest, entites.Error{
			Error: err.Error(),
		})
		return
	}

	newIP := c.Request.RemoteAddr

	newAccess, err := h.service.Refresh(c.Request.Context(), userID, req.RefreshToken, newIP)
	if err != nil {
		log.Errorw("error with refreshing tokens", zap.Error(err))
		c.JSON(http.StatusBadRequest, entites.Error{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entites.Response{
		AccessToken:  newAccess,
		RefreshToken: req.RefreshToken,
	})
}
