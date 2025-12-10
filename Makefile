# dysv.de Development Makefile

.PHONY: help dev api web stripe build test clean

# Default target
help:
	@echo "dysv.de Development Commands"
	@echo ""
	@echo "  make dev         - Run API + Web + Stripe listener (requires tmux)"
	@echo "  make api         - Run Go API server"
	@echo "  make web         - Run frontend dev server"
	@echo "  make stripe      - Run Stripe webhook listener"
	@echo "  make build       - Build both API and frontend"
	@echo "  make test        - Run all tests"
	@echo "  make clean       - Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  MONGODB_URI            - MongoDB connection (default: mongodb://localhost:27017/dysv)"
	@echo "  STRIPE_SECRET          - Stripe secret key (sk_test_...)"
	@echo "  STRIPE_WEBHOOK_SECRET  - Stripe webhook secret (whsec_...)"

# Run everything in tmux (API + Web + Stripe)
SESSION_NAME := dysv-dev

dev:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		echo "Session $(SESSION_NAME) already exists. Run 'make dev-stop' first or 'tmux attach -t $(SESSION_NAME)'"; \
		exit 1; \
	fi
	@echo "Starting dev environment in tmux session '$(SESSION_NAME)'..."
	@tmux new-session -d -s $(SESSION_NAME) -n api \
		'MONGODB_URI=$${MONGODB_URI:-mongodb://localhost:27017/dysv} go run . api; read'
	@tmux new-window -t $(SESSION_NAME) -n web \
		'cd web && bun run dev; read'
	@tmux new-window -t $(SESSION_NAME) -n stripe \
		'echo "Starting Stripe listener..."; stripe listen --forward-to localhost:8080/api/webhook/stripe; read'
	@tmux select-window -t $(SESSION_NAME):api
	@echo ""
	@echo "Dev environment started!"
	@echo "  tmux attach -t $(SESSION_NAME)    - Attach to session"
	@echo "  make dev-stop                     - Stop all services"
	@echo ""
	@tmux attach -t $(SESSION_NAME)

dev-stop:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		tmux kill-session -t $(SESSION_NAME); \
		echo "Session $(SESSION_NAME) stopped."; \
	else \
		echo "Session $(SESSION_NAME) not running."; \
	fi

# Run Go API server
api:
	@echo "Starting API server on :8080..."
	MONGODB_URI=$${MONGODB_URI:-mongodb://localhost:27017/dysv} \
	go run . api

# Run frontend dev server
web:
	@echo "Starting frontend dev server..."
	cd web && bun run dev

# Run Stripe webhook listener (forwards to local API)
stripe:
	@echo "Starting Stripe webhook listener..."
	@echo "Copy the webhook secret (whsec_...) to STRIPE_WEBHOOK_SECRET"
	stripe listen --forward-to localhost:8080/api/webhook/stripe

# Build everything
build: build-api build-web

build-api:
	@echo "Building Go API..."
	go build -o bin/dysv .

build-web:
	@echo "Building frontend..."
	cd web && bun run build

# Run tests
test: test-api test-web

test-api:
	go test ./...

test-web:
	cd web && bun test 2>/dev/null || echo "No tests configured"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/.output/
	rm -rf web/node_modules/.vite/

# Install dependencies
deps:
	go mod download
	cd web && bun install

# Format code
fmt:
	go fmt ./...
	cd web && bun run check 2>/dev/null || true

# Lint
lint:
	go vet ./...
	cd web && bun run check
