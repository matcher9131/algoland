# algoland — Next.js (web) + Go (api) deployed as one Vercel project.

GO ?= $(HOME)/.local/go/bin/go

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install all dependencies
	npm install
	cd web && npm install
	cd api && $(GO) mod download

.PHONY: dev
dev: ## Run every service together via the Vercel CLI
	npx vercel dev

.PHONY: dev-web
dev-web: ## Run only the Next.js frontend
	cd web && npm run dev

.PHONY: dev-api
dev-api: ## Run only the Go backend
	cd api && $(GO) run ./cmd/server

.PHONY: db-start
db-start: ## Start the local Supabase stack
	npx supabase start

.PHONY: db-stop
db-stop: ## Stop the local Supabase stack
	npx supabase stop

.PHONY: db-reset
db-reset: ## Recreate the local database from migrations
	npx supabase db reset

.PHONY: db-push
db-push: ## Apply migrations to the linked remote project
	npx supabase db push

.PHONY: lint
lint: ## Lint both services
	cd web && npm run lint
	cd api && $(GO) vet ./...

.PHONY: typecheck
typecheck: ## Type-check the frontend
	cd web && npm run typecheck

.PHONY: test
test: ## Run the Go unit tests
	cd api && $(GO) test ./...

.PHONY: test-db
test-db: ## Run the Go tests including the ones needing a live database
	@test -n "$(ALGOLAND_TEST_DATABASE_URL)" || { \
		echo "ALGOLAND_TEST_DATABASE_URL is required."; \
		echo "Get it from: npx supabase status -o env | grep DB_URL"; \
		exit 1; }
	cd api && $(GO) test ./... -v

.PHONY: build
build: ## Build both services
	cd web && npm run build
	cd api && $(GO) build -o bin/server ./cmd/server

.PHONY: docker-build
docker-build: ## Build the API container image exactly as Vercel does
	docker build -f api/Dockerfile.vercel -t algoland-api:local api

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf web/.next api/bin
