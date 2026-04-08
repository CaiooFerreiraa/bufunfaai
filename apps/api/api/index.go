package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/bootstrap"
)

var (
	appOnce sync.Once
	app     *bootstrap.App
	appErr  error
)

func Handler(writer http.ResponseWriter, request *http.Request) {
	appOnce.Do(func() {
		app, appErr = bootstrap.Build(context.Background())
	})

	if appErr != nil {
		http.Error(writer, "failed to bootstrap api", http.StatusInternalServerError)
		return
	}

	app.Server.Engine().ServeHTTP(writer, request)
}
