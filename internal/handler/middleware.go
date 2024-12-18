package handler

import (
	"awesomeProject/internal/entites"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusBadRequest, entites.Error{Error: "Missing token"})
		return
	}
	splittedHeader := strings.Split(tokenString, " ")
	if len(splittedHeader) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusBadRequest, entites.Error{Error: "Missing token"})
		return
	}
	c.Set("accessToken", splittedHeader[1])
}
