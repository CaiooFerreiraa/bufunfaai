package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/health/interfaces/http/handler"
)

func Register(engine *gin.Engine, handler *handler.Handler) {
	engine.GET("/health", handler.Health)
	engine.GET("/ready", handler.Ready)
}
