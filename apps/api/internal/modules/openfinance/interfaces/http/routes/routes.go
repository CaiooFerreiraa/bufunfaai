package routes

import (
	"github.com/gin-gonic/gin"

	ofhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/handler"
)

func Register(publicGroup *gin.RouterGroup, protectedGroup *gin.RouterGroup, handler *ofhandler.Handler) {
	publicOpenFinanceGroup := publicGroup.Group("/open-finance")
	publicOpenFinanceGroup.GET("/callback", handler.Callback)
	publicOpenFinanceGroup.POST("/callback", handler.Callback)

	openFinanceGroup := protectedGroup.Group("/open-finance")
	openFinanceGroup.GET("/institutions", handler.ListInstitutions)
	openFinanceGroup.GET("/institutions/:id", handler.GetInstitution)
	openFinanceGroup.POST("/consents", handler.CreateConsent)
	openFinanceGroup.GET("/consents/:id", handler.GetConsent)
	openFinanceGroup.POST("/consents/:id/authorize", handler.AuthorizeConsent)
	openFinanceGroup.POST("/consents/:id/connect-token", handler.CreateConnectToken)
	openFinanceGroup.POST("/consents/:id/complete", handler.CompleteConsent)
	openFinanceGroup.POST("/consents/:id/revoke", handler.RevokeConsent)
	openFinanceGroup.GET("/connections", handler.ListConnections)
	openFinanceGroup.GET("/connections/:id", handler.GetConnection)
	openFinanceGroup.DELETE("/connections/:id", handler.DeleteConnection)
	openFinanceGroup.POST("/connections/:id/sync", handler.SyncConnection)
	openFinanceGroup.GET("/connections/:id/sync-status", handler.SyncStatus)
}

func RegisterInternal(internalGroup *gin.RouterGroup, handler *ofhandler.Handler) {
	internalOpenFinanceGroup := internalGroup.Group("/open-finance")
	internalOpenFinanceGroup.POST("/reconcile", handler.ReconcileConnections)
}
