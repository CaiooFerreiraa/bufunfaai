package routes

import (
	"github.com/gin-gonic/gin"

	analyticshandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/handler"
)

func Register(protectedGroup *gin.RouterGroup, handler *analyticshandler.Handler) {
	analyticsGroup := protectedGroup.Group("/analytics")
	analyticsGroup.GET("/overview", handler.Overview)
	analyticsGroup.GET("/budgets", handler.Budgets)
	analyticsGroup.GET("/score/latest", handler.LatestScore)
	analyticsGroup.GET("/forecast", handler.Forecast)
	analyticsGroup.GET("/insights", handler.Insights)
	analyticsGroup.GET("/anomalies", handler.Anomalies)
	analyticsGroup.GET("/goals", handler.Goals)
	analyticsGroup.POST("/goals", handler.CreateGoal)
	analyticsGroup.PATCH("/goals/:id", handler.UpdateGoal)
}
