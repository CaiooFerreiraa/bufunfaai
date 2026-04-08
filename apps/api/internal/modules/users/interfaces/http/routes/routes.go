package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/handler"
)

func Register(group *gin.RouterGroup, handler *handler.Handler) {
	usersGroup := group.Group("/users")
	usersGroup.GET("/me", handler.Me)
	usersGroup.PATCH("/me", handler.UpdateMe)
}
