package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey string = "request_id"

func RequestID() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID := context.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		context.Set(requestIDKey, requestID)
		context.Writer.Header().Set("X-Request-ID", requestID)
		context.Next()
	}
}
