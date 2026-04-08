package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/consents/interfaces/http/handler"
)

func Register(group *gin.RouterGroup, handler *handler.Handler) {
	consentsGroup := group.Group("/consents")
	consentsGroup.GET("", handler.List)
	consentsGroup.POST("", handler.Create)
	consentsGroup.DELETE("/:id", handler.Delete)
}
