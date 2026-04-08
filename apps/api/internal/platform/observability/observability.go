package observability

import (
	"context"
	"time"

	ddgin "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin/v2"
	ddtracer "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/logger"
)

type ShutdownFunc func(context.Context)

func Initialize(cfg config.Config, appLogger *logger.Logger) (ShutdownFunc, error) {
	shutdowns := make([]ShutdownFunc, 0, 2)

	if cfg.SentryDSN != "" {
		options := sentry.ClientOptions{
			Dsn:         cfg.SentryDSN,
			Environment: cfg.SentryEnvironment,
			Release:     cfg.SentryRelease,
		}

		if cfg.SentryTracesSampleRate > 0 {
			options.EnableTracing = true
			options.TracesSampleRate = cfg.SentryTracesSampleRate
		}

		if err := sentry.Init(options); err != nil {
			return nil, err
		}

		appLogger.Info("sentry initialized", "environment", cfg.SentryEnvironment)
		shutdowns = append(shutdowns, func(_ context.Context) {
			sentry.Flush(2 * time.Second)
		})
	}

	if cfg.DatadogTraceEnabled {
		options := []ddtracer.StartOption{
			ddtracer.WithEnv(cfg.DatadogEnv),
			ddtracer.WithService(cfg.DatadogService),
		}

		if cfg.DatadogVersion != "" {
			options = append(options, ddtracer.WithServiceVersion(cfg.DatadogVersion))
		}

		ddtracer.Start(options...)
		appLogger.Info("datadog tracer initialized", "service", cfg.DatadogService, "env", cfg.DatadogEnv)
		shutdowns = append(shutdowns, func(_ context.Context) {
			ddtracer.Stop()
		})
	}

	return func(ctx context.Context) {
		for index := len(shutdowns) - 1; index >= 0; index-- {
			shutdowns[index](ctx)
		}
	}, nil
}

func AttachGin(engine *gin.Engine, cfg config.Config) {
	if cfg.SentryDSN != "" {
		engine.Use(sentrygin.New(sentrygin.Options{}))
	}

	if cfg.DatadogTraceEnabled {
		engine.Use(ddgin.Middleware(cfg.DatadogService))
	}
}
