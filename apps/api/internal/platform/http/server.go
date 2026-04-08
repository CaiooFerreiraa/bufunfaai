package http

import (
	"context"
	"fmt"
	stdhttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/logger"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
)

type Server struct {
	engine *gin.Engine
	server *stdhttp.Server
}

func NewServer(cfg config.Config, appLogger *logger.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(appLogger))

	return &Server{
		engine: engine,
		server: &stdhttp.Server{
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           engine,
			ReadHeaderTimeout: cfg.HTTPReadTimeout,
			WriteTimeout:      cfg.HTTPWriteTimeout,
		},
	}
}

func (server *Server) Engine() *gin.Engine {
	return server.engine
}

func (server *Server) Run() error {
	return server.server.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}
