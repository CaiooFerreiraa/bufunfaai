package handler

import (
	stdcontext "context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type ReadinessProbe interface {
	Name() string
	Check(ctx stdcontext.Context) error
}

type Handler struct {
	config config.Config
	probes []ReadinessProbe
}

func NewHandler(cfg config.Config, probes []ReadinessProbe) *Handler {
	return &Handler{
		config: cfg,
		probes: probes,
	}
}

func (handler *Handler) Health(ginContext *gin.Context) {
	response.OK(ginContext, gin.H{
		"status":      "ok",
		"service":     handler.config.AppName,
		"version":     "0.1.0",
		"environment": handler.config.Environment,
	})
}

func (handler *Handler) Ready(ginContext *gin.Context) {
	checkContext, cancel := stdcontext.WithTimeout(ginContext.Request.Context(), 2*time.Second)
	defer cancel()

	failures := make([]gin.H, 0)
	for _, probe := range handler.probes {
		if err := probe.Check(checkContext); err != nil {
			failures = append(failures, gin.H{
				"name":   probe.Name(),
				"status": "down",
			})
		}
	}

	if len(failures) > 0 {
		response.Success(ginContext, http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"services": failures,
		})
		return
	}

	response.OK(ginContext, gin.H{
		"status": "ready",
	})
}
