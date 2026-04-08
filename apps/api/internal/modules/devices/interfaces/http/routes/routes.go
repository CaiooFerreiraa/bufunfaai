package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/handler"
)

func Register(group *gin.RouterGroup, handler *handler.Handler) {
	devicesGroup := group.Group("/devices")
	devicesGroup.GET("", handler.List)
	devicesGroup.DELETE("/:id", handler.Delete)
}
