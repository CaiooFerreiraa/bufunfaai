package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/bootstrap"
)

func main() {
	ctx := context.Background()
	app, err := bootstrap.Build(ctx)
	if err != nil {
		log.Fatalf("failed to bootstrap worker: %v", err)
	}
	defer app.Close()

	jobName := strings.TrimSpace(os.Getenv("WORKER_JOB"))
	if jobName == "" {
		jobName = "openfinance-reconcile"
	}

	switch jobName {
	case "openfinance-reconcile":
		result, appError := app.OpenFinanceUseCases.ReconcileConnections(ctx, app.Config.WorkerBatchSize)
		if appError != nil {
			log.Fatalf("worker job failed: %v", appError)
		}

		app.Logger.Info(
			"worker job completed",
			"job", jobName,
			"processed", result.Processed,
			"successful", result.Successful,
			"failed", result.Failed,
			"jobs_created", result.JobsCreated,
		)
	default:
		log.Fatalf("unknown worker job: %s", jobName)
	}
}
