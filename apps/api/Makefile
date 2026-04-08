APP_NAME=finance-api

.PHONY: dev run build worker reconcile test lint fmt tidy up down logs

dev:
	go run ./cmd/api

run:
	go run ./cmd/api

build:
	go build ./cmd/api

worker:
	go run ./cmd/worker

reconcile:
	$$env:WORKER_JOB='openfinance-reconcile'; go run ./cmd/worker

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w ./cmd ./internal ./tests

tidy:
	go mod tidy

up:
	docker compose up -d postgres redis

down:
	docker compose down

logs:
	docker compose logs -f api
