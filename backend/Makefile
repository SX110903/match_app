.PHONY: dev build test test-race migrate migrate-down generate-keys lint docker-build docker-push seed help

BINARY=server
BUILD_DIR=./bin
CMD=./cmd/server

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Start with hot reload (requires Docker)
	docker compose -f deployments/docker-compose.dev.yml up

build: ## Compile optimized binary
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY) $(CMD)

test: ## Run all tests with coverage
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-race: ## Run tests with race detector
	go test -race ./...

migrate: ## Apply all pending migrations
	@source .env && bash scripts/migrate.sh

generate-keys: ## Generate RSA keys for JWT and AES key for encryption
	bash scripts/generate-keys.sh

lint: ## Run golangci-lint
	golangci-lint run ./...

docker-build: ## Build production Docker image
	docker build -f deployments/docker/Dockerfile -t matchhub-backend:latest .

docker-push: ## Push image to registry (set REGISTRY env var)
	docker push $(REGISTRY)/matchhub-backend:latest

seed: ## Populate DB with test data
	bash scripts/seed.sh

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) coverage.out coverage.html
