# algoland — Next.js (web) + Go (api) deployed as one Vercel project.

GO ?= $(HOME)/.local/go/bin/go

# Pinned; run via `go run` so no separate install is needed.
GOLANGCI_LINT_VERSION ?= v2.15.0
GOLANGCI_LINT ?= $(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# Chromium needs a handful of shared libraries that `playwright install-deps`
# would install with apt and root. These paths hold a rootless copy instead;
# playwright.config.ts picks them up automatically when present.
E2E_DEPS_DIR ?= $(HOME)/.local/playwright-deps
E2E_DEB_DIR ?= $(HOME)/.cache/algoland-debs

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) \
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

.PHONY: fmt
fmt: fmt-web fmt-api ## Format both services

.PHONY: fmt-web
fmt-web: ## Format web (Prettier)
	cd web && npm run format

.PHONY: fmt-check-web
fmt-check-web: ## Check web formatting (Prettier)
	cd web && npm run format:check

.PHONY: fmt-api
fmt-api: ## Format api (gofmt, goimports)
	cd api && $(GOLANGCI_LINT) fmt ./...

.PHONY: lint
lint: lint-web lint-api ## Lint both services

.PHONY: lint-web
lint-web: ## Lint web
	cd web && npm run lint

.PHONY: lint-api
lint-api: ## Lint api
	cd api && $(GOLANGCI_LINT) run ./...

.PHONY: typecheck
typecheck: ## Type-check the frontend
	cd web && npm run typecheck

.PHONY: test
test: test-api test-web ## Run every test except E2E

.PHONY: test-api
test-api: ## Run the Go tests (starts a Postgres container)
	cd api && $(GO) test ./...

.PHONY: test-api-short
test-api-short: ## Run only the Go tests that need no Docker
	cd api && $(GO) test -short ./...

.PHONY: test-web
test-web: ## Run the Vitest suites
	cd web && npm test

.PHONY: e2e-deps
e2e-deps: ## Install Chromium's system libraries into $$HOME (no root needed)
	@mkdir -p $(E2E_DEB_DIR) $(E2E_DEPS_DIR)
	cd $(E2E_DEB_DIR) && apt-get download libnss3 libnspr4 libasound2t64
	cd $(E2E_DEB_DIR) && for d in *.deb; do dpkg -x "$$d" $(E2E_DEPS_DIR); done
	cd web && npx playwright install chromium
	@echo "Chromium dependencies unpacked into $(E2E_DEPS_DIR)."

.PHONY: test-e2e
test-e2e: ## Run the Playwright E2E suite against the full local stack
	@test -n "$(DATABASE_URL)" || { \
		echo "DATABASE_URL and SUPABASE_URL are required for the backend specs."; \
		echo "Load them with: export $$(npx supabase status -o env | tr -d '\"' | xargs)"; \
		exit 1; }
	cd web && npm run test:e2e

.PHONY: coverage
coverage: ## Report coverage for both services
	cd api && $(GO) test ./... -coverprofile=coverage.out -coverpkg=./... \
		&& $(GO) tool cover -func=coverage.out | tail -1
	cd web && npm run test:coverage

.PHONY: build
build: ## Build both services
	cd web && npm run build
	cd api && $(GO) build -o bin/server ./cmd/server

.PHONY: docker-build
docker-build: ## Build the API container image exactly as Vercel does
	docker build -f api/Dockerfile.vercel -t algoland-api:local api

.PHONY: clean
clean: ## Remove build artifacts and test output
	rm -rf web/.next api/bin api/coverage.out
	rm -rf web/coverage web/playwright-report web/test-results
