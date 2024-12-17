package service

import (
	"awesomeProject/internal/logger"
	"awesomeProject/internal/store"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret string
	store     store.Auth
}

func NewAuthService(store store.Auth, jwtSecret string) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) GenerateTokens(ctx context.Context, userId string, clientIp string) (string, string, error) {
	log := logger.LoggerFromContext(ctx)
	user, err := s.store.GetUserById(ctx, userId)
	if err != nil {
		log.Errorw("error with getting user by id", zap.Error(err))
		return "", "", err
	}

	if user == nil {
		log.Errorw("there is no user with this id", zap.Error(err))
		return "", "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":    userId,
		"client_ip":  clientIp,
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().Add(15 * time.Minute).Unix(), //todo из конфига
	})

	jwtToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Errorw("error with generating jwt token", zap.Error(err))
		return "", "", err
	}
	rawRefresh := make([]byte, 32)
	_, err = rand.Read(rawRefresh)
	if err != nil {
		log.Errorw("error with generating refresh token", zap.Error(err))
		return "", "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(rawRefresh)
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Errorw("error with generating refresh token", zap.Error(err))
		return "", "", err
	}

	if err := s.store.InsertUserInfo(ctx, userId, clientIp, string(hash)); err != nil {
		log.Errorw("error with inserting user info", zap.Error(err))
		return "", "", err
	}

	return jwtToken, string(hash), nil
}

func (s *AuthService) ExtractUserIDFromAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", jwt.ErrInvalidKey
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, "", errors.New("user_id not found in token")
	}

	return claims, userID, nil
}

func (as *AuthService) Refresh(ctx context.Context, userID string, providedRefreshToken string, newClientIP string) (newAccess string, err error) {
	rt, err := as.store.GetUserById(ctx, userID)
	if err != nil {
		return "", err
	}
	if rt == nil {
		return "", errors.New("no refresh token found")
	}

	// Проверяем хэш //todo тут хэш разный для рефреша и jwt
	err = bcrypt.CompareHashAndPassword([]byte(rt.Hash), []byte(providedRefreshToken))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Проверяем срок действия
	if time.Now().After(rt.) { //todo тут нужно получать рефреш для конкретного пользователя
		return "", errors.New("refresh token expired")
	}

	// Проверяем IP (если нужно). Если IP изменился — высылаем предупреждение (опционально)
	if rt.ClientIP != newClientIP {
		user, _ := as.store.GetUserById(ctx, userID)
		if user != nil {
			as.EmailService.SendWarning(user.Email, rt.ClientIP, newClientIP)
		}
	}

	// Генерируем новый Access Token
	newAccess, err = security.GenerateAccessToken(userID, newClientIP)
	if err != nil {
		return "", err
	}

	// Refresh token остается тем же
	return newAccess, nil
}
