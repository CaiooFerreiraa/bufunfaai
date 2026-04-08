package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/logger"
)

func RequestLogger(appLogger *logger.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		startedAt := time.Now()

		context.Next()

		appLogger.Info(
			"http request",
			"request_id", requestIDFromContext(context),
			"method", context.Request.Method,
			"path", context.Request.URL.Path,
			"status", context.Writer.Status(),
			"latency_ms", time.Since(startedAt).Milliseconds(),
			"client_ip", context.ClientIP(),
			"user_agent", context.Request.UserAgent(),
		)
	}
}

func requestIDFromContext(context *gin.Context) string {
	value, ok := context.Get(requestIDKey)
	if !ok {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
