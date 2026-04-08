package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireInternalSecret(secret string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if secret == "" {
			context.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_SECRET_NOT_CONFIGURED",
					"message": "segredo interno nao configurado",
				},
			})
			return
		}

		providedSecret := context.GetHeader("X-Internal-Secret")
		if providedSecret == "" {
			providedSecret = context.GetHeader("Authorization")
			if len(providedSecret) > 7 && providedSecret[:7] == "Bearer " {
				providedSecret = providedSecret[7:]
			}
		}

		if providedSecret != secret {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_INTERNAL_SECRET",
					"message": "autorizacao interna invalida",
				},
			})
			return
		}

		context.Next()
	}
}
