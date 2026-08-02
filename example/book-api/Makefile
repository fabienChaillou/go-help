.PHONY: up down setup run test test-integration lint seed

up:
	@docker run -d \
  --name book-api-postgres \
  -e POSTGRES_USER=bookuser \
  -e POSTGRES_PASSWORD=bookpass \
  -e POSTGRES_DB=bookdb \
  -p 5432:5432 \
  postgres:16-alpine

down:
	docker compose down

setup:
	go mod tidy

run: up
	go run ./cmd/server/main.go

test:
	go test ./internal/... -v -count=1

test-integration:
	go test ./test/integration/... -v -count=1 -timeout 120s

lint:
	golangci-lint run ./...

seed: up
	go run ./cmd/seed/main.go