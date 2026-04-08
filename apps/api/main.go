package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/bootstrap"
)

func main() {
	ctx := context.Background()
	app, err := bootstrap.Build(ctx)
	if err != nil {
		log.Fatalf("failed to bootstrap api: %v", err)
	}
	defer app.Close()

	app.Logger.Info("starting api", "port", app.Config.Port, "env", app.Config.Environment)

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownSignals
		shutdownContext, cancel := context.WithTimeout(context.Background(), app.Config.HTTPShutdownTimeout)
		defer cancel()

		if err := app.Server.Shutdown(shutdownContext); err != nil {
			app.Logger.Error("failed to shutdown server", "error", err.Error())
		}
	}()

	if err := app.Server.Run(); err != nil {
		app.Logger.Error("server stopped", "error", err.Error())
	}
}
