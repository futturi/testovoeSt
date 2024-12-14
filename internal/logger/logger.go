package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.SugaredLogger); ok {
		return l
	}
	return zap.S()
}

func LoggerMiddleware(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.Clone(context.WithValue(c.Request.Context(), ctxLogger{}, log))
		c.Next()
	}
}

func InitLogger() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	logger, _ := cfg.Build()
	return logger.Sugar()
}
