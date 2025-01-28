# Environment variables
SHELL := /bin/bash

# Frontend Applications
APPS_DIR := apps
APPS := $(shell find $(APPS_DIR) -maxdepth 1 -type d -not -path '*/\.*' -not -path '$(APPS_DIR)' -exec basename {} \;)

# Backend services
BACKENDS_DIR := backends
SERVICES_WITHOUT_LIVE_TEST := api
SERVICES_WITH_LIVE_TEST := $(filter-out $(SERVICES_WITHOUT_LIVE_TEST),$(shell find $(BACKENDS_DIR) -maxdepth 1 -type d -not -path '*/\.*' -not -path '$(BACKENDS_DIR)' -not -path '$(BACKENDS_DIR)/internal' -not -path '$(BACKENDS_DIR)/internal' -exec basename {} \;))
SERVICES := $(SERVICES_WITHOUT_LIVE_TEST) $(SERVICES_WITH_LIVE_TEST)

# Temporary build directory
BUILD_DIR := ./tmp

# Build flags
GO_BUILD_FLAGS := -v  # Verbose Go build output
GO_TEST_FLAGS := -race -v  # Go test flags with race detection

# CI-specific configurations
ifeq ($(CI),true)
	GOLANGCI_LINT :=
	BIOME_FLAGS := --diagnostic-level=warn
else
	GOLANGCI_LINT := golangci-lint run
	BIOME_FLAGS := --diagnostic-level=warn --write --unsafe
endif

.PHONY: all check check-be check-fe clean db-migrate dev dev-be dev-fe domains down gen gen-be gen-fe help preview reset setup test test-be up

# Default target - sets up environment and starts development
all: setup dev

# Run all code checks
check: gen
	@pnpm --silent concurrently -g -r \
		"make check-be" \
		"make check-fe"

# Run backend code checks (buf, go fmt, go vet, golangci-lint)
check-be:
	@buf format -w
	@buf lint
	@go fmt ./...
	@go vet ./...
	@$(GOLANGCI_LINT)

# Run frontend code checks using biome
check-fe:
	@pnpm --silent biome check $(BIOME_FLAGS) .
	@for app in $(APPS); do pnpm --silent --filter=./$(APPS_DIR)/$$app tsc; done

# Clean build artifacts and temporary files
clean:
	@rm -rf tmp **/**/{.react-router,build,tsconfig.tsbuildinfo}

# Run database migrations for all services
db-migrate:
	@echo "🔄 Running Database Migrations"
	@for service in $(SERVICES); do \
		echo ""; \
		echo "──────────────────────────────────────────────────────────────────────────────"; \
		echo "📦 Migrating '$$service' database..."; \
		echo "──────────────────────────────────────────────────────────────────────────────"; \
		go run ./$(BACKENDS_DIR)/$$service db:migrate; \
	done
	@echo "✅ All migrations completed"

# Create template databases for all services
db-templates:
	@echo "🔄 Creating template databases..."
	@for db in $(SERVICES_WITHOUT_LIVE_TEST); do docker compose exec -it db psql -U postgres -q -c "CREATE DATABASE template_$$db WITH TEMPLATE '$$db';"; done
	@for db in $(SERVICES_WITH_LIVE_TEST); do docker compose exec -it db psql -U postgres -q -c "CREATE DATABASE template_$$db WITH TEMPLATE '$${db}_live';"; done
	@echo "✅ All template databases created"

# Start development environment with all services
dev: gen
	@mkdir -p $(BUILD_DIR)
	@pnpm --silent concurrently -g -r \
		"pnpm --silent chokidar './packages/api/src/contracts/*.json' -c 'pnpm --silent --filter=./packages/api gen' --silent" \
		"pnpm --silent chokidar '$(BACKENDS_DIR)/api/main.go' '$(BACKENDS_DIR)/api/internal/handler/**/*.go' -c 'go run ./$(BACKENDS_DIR)/api gen:openapi' --silent" \
		"pnpm --silent chokidar '$(BACKENDS_DIR)/internal/pb/**/*.proto' -c 'make gen-be' --silent" \
		"pnpm --silent --filter=./packages/ui storybook" \
		$(foreach app,$(APPS),"make dev-fe APP_SERVICE=$(app)") \
		$(foreach service,$(SERVICES),"make dev-be APP_SERVICE=$(service)")

# Run backend service with hot reload using air
dev-be:
	@air -build.bin='$(BUILD_DIR)/$(APP_SERVICE) start --worker' \
		-build.cmd='go build -o $(BUILD_DIR)/$(APP_SERVICE) ./$(BACKENDS_DIR)/$(APP_SERVICE)' \
		-build.delay=350 \
		-build.exclude_dir=$(BACKENDS_DIR)/internal/core/testdata,$(BACKENDS_DIR)/internal/pb \
		-build.include_ext='css,go,html,js,json,sql,toml,tpl,tmpl,yaml,yml' \
		-build.include_dir=$(BACKENDS_DIR)/$(APP_SERVICE),$(BACKENDS_DIR)/internal \
		-log.main_only=true

# Run frontend app in development mode
dev-fe:
	@pnpm --silent --filter=./$(APPS_DIR)/$(APP_SERVICE) dev

# Show all local domains
domains:
	@echo "🌐 Local Domains"
	@echo "─────────────────────────────────────────────────"
	@echo "📱 Frontend Apps"
	@echo "   • Storybook:        http://localhost:2995"
	@echo "   • Dashboard:        http://localhost:3000"
	@echo "   • Mailer Preview:   http://localhost:3001/mailer/preview"
	@echo ""
	@echo "🔧 Backend Services"
	@echo "   • API:              http://localhost:3001"
	@echo "   • Payment:          http://localhost:3002"
	@echo ""
	@echo "🛠️  Infrastructure"
	@echo "   • Postgres:         localhost:5432"
	@echo "   • Redis:            localhost:6379"
	@echo "   • Kafka:            localhost:9092"
	@echo "   • Kafka UI:         http://localhost:8080"
	@echo "   • Mailpit:          http://localhost:8025"
	@echo "   • MinIO:            http://localhost:9000"
	@echo "   • MinIO Admin:      http://localhost:9001"
	@echo "─────────────────────────────────────────────────"

# Stop and remove infrastructure containers
down:
	@docker compose --profile=infra down --remove-orphans --timeout 0

# Generate all code (protobuf, OpenAPI, and API clients)
gen:
	@make gen-be
	@make gen-fe

# Generate backend code and contracts
gen-be:
	@buf generate --clean
	@go run ./$(BACKENDS_DIR)/api gen:openapi
	@pnpm --silent --filter=./packages/api gen

# Generate frontend code and types
gen-fe:
	@for app in $(APPS); do	pnpm --silent --filter=./$(APPS_DIR)/$$app react-router typegen; done

# Show available make commands
help:
	@echo "🛠️  Available Make Commands"
	@echo "────────────────────────────────────────────────────────────────────────────────"
	@echo "Development Commands:"
	@echo "   • make all              - Set up environment and start development"
	@echo "   • make dev              - Start development environment with all services"
	@echo "   • make dev-be           - Run backend service with hot reload"
	@echo "   • make dev-fe           - Run frontend app in development mode"
	@echo "   • make domains          - Show all local domains and ports"
	@echo ""
	@echo "Infrastructure Commands:"
	@echo "   • make up               - Start infrastructure containers and run migrations"
	@echo "   • make down             - Stop and remove infrastructure containers"
	@echo "   • make reset            - Reset infrastructure (down + up)"
	@echo ""
	@echo "Code Generation Commands:"
	@echo "   • make gen              - Generate all code (protobuf, OpenAPI, API clients)"
	@echo "   • make gen-be           - Generate backend code and contracts"
	@echo "   • make gen-fe           - Generate frontend code and types"
	@echo ""
	@echo "Testing Commands:"
	@echo "   • make test             - Run all tests"
	@echo "   • make test-be          - Run backend tests with race detection"
	@echo ""
	@echo "Code Quality Commands:"
	@echo "   • make check            - Run all code checks"
	@echo "   • make check-be         - Run backend code checks"
	@echo "   • make check-fe         - Run frontend code checks"
	@echo ""
	@echo "Database Commands:"
	@echo "   • make db-migrate       - Run database migrations for all services"
	@echo ""
	@echo "Other Commands:"
	@echo "   • make setup            - Install dependencies and setup pre-commit hooks"
	@echo "   • make clean            - Clean build artifacts and temporary files"
	@echo "   • make preview          - Preview application in Docker"
	@echo "────────────────────────────────────────────────────────────────────────────────"
	@echo "💡 Tip: Use 'make <command>' to run a command"

# Preview application in Docker
preview:
	@docker compose --profile=app up --build
	@docker compose --profile=app down

# Reset infrastructure (down + up)
reset: down up

# Install dependencies and setup pre-commit hooks
setup:
	@pre-commit install --hook-type pre-push
	@go mod download
	@pnpm i

# Run all tests
test: test-be

# Run backend tests with race detection
test-be:
	@gotestsum --format=short -- $(GO_TEST_FLAGS) ./...

# Start infrastructure containers and run migrations
up:
	@docker compose --profile=infra up --wait
	@make db-migrate
	@make db-templates

# Run individual services directly
$(SERVICES):
	@go run ./$(BACKENDS_DIR)/$@ $(filter-out $@ --,$(MAKECMDGOALS))

%:
	@:
