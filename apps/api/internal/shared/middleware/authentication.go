package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	platformauth "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/auth"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

const (
	userIDContextKey    string = "user_id"
	sessionIDContextKey string = "session_id"
)

func RequireAuthentication(jwtService *platformauth.JWTService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader("Authorization")
		if authorizationHeader == "" {
			response.Error(context, sharederrors.New("UNAUTHORIZED", "autenticacao obrigatoria", http.StatusUnauthorized))
			context.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(context, sharederrors.New("UNAUTHORIZED", "token invalido", http.StatusUnauthorized))
			context.Abort()
			return
		}

		claims, err := jwtService.Parse(parts[1])
		if err != nil {
			response.Error(context, sharederrors.New("UNAUTHORIZED", "token invalido ou expirado", http.StatusUnauthorized))
			context.Abort()
			return
		}

		context.Set(userIDContextKey, claims.Subject)
		context.Set(sessionIDContextKey, claims.SessionID)
		context.Next()
	}
}

func CurrentUserID(context *gin.Context) string {
	value, ok := context.Get(userIDContextKey)
	if !ok {
		return ""
	}

	userID, ok := value.(string)
	if !ok {
		return ""
	}

	return userID
}
