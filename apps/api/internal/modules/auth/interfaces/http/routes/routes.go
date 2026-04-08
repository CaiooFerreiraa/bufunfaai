package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/handler"
)

func Register(publicGroup *gin.RouterGroup, protectedGroup *gin.RouterGroup, handler *handler.Handler) {
	authGroup := publicGroup.Group("/auth")
	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/refresh", handler.Refresh)
	authGroup.POST("/logout", handler.Logout)

	protectedAuthGroup := protectedGroup.Group("/auth")
	protectedAuthGroup.POST("/logout-all", handler.LogoutAll)
}
