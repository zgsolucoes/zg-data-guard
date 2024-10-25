include .env
# Functions
urlencode = $(shell python3 -c "import urllib.parse; print(urllib.parse.quote('$(1)'))")

# Variables
APP_NAME=zg-data-guard
BUILD_FILE = build.properties
BUILD_DIR=dist
TEST_REPORTS_DIR=testdata/reports
DOCS_DIR=docs
MIN_COVERAGE=45
MIN_CORE_COVERAGE=97
# Get Git branch name
BUILD_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
# Get latest commit hash
BUILD_HASH = $(shell git rev-parse --short HEAD)
# Define build version with priority for "version" variable
buildVersion = $${version:-$(BUILD_BRANCH)-$(BUILD_HASH)}

# Tasks
default: run

create_migration:
	@echo "==> Creating new migration..."
	migrate create -ext=sql -dir=internal/database/migrations -seq zg_data_guard

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://$(call urlencode,$(DATABASE_USER)):$(call urlencode,$(DATABASE_PASSWORD))@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://$(call urlencode,$(DATABASE_USER)):$(call urlencode,$(DATABASE_PASSWORD))@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down 1

migrate_force:
	migrate -path=internal/database/migrations -database "postgresql://$(call urlencode,$(DATABASE_USER)):$(call urlencode,$(DATABASE_PASSWORD))@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose force $(version)

.PHONY: install
install:
	@echo "==> Installing application dependencies..."
	@go mod tidy

run:
	@echo "==> Running the application..."
	@echo "version=$(buildVersion)" > $(BUILD_FILE)
	@echo "buildTime=$(shell date -u '+%Y-%m-%d %H:%M:%S')" >> $(BUILD_FILE)
	@go run cmd/zg-data-guard/main.go

run-with-docs: docs run

docs:
	@echo "==> Removing old API documentation..."
	@rm -rf $(DOCS_DIR)
	@echo "> Generating new swagger API documentation..."
	@swag init -g cmd/zg-data-guard/main.go

.PHONY: lint
lint:
	@echo "==> Running linter..."
	@golangci-lint run
	@echo "> Linter completed successfully!"

test:
	@echo "==> Running tests..."
	@go test ./...

.PHONY: test-verbose
test-verbose:
	@echo "==> Running tests with verbose..."
	@go test -v ./...

.PHONY: html-coverage
html-coverage:
	@echo "==> Running tests with coverage and generating HTML report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY: test-count
test-count:
	@echo "==> Running tests without cache and with counting..."
	@go test ./... -count=1 -v | grep -c RUN

.PHONY: coverage
coverage:
	@echo "==> Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@echo "> Generating HTML report..."
	@mkdir -p $(TEST_REPORTS_DIR)
	@go tool cover -html=coverage.out -o $(TEST_REPORTS_DIR)/coverage.html
	@echo "> Checking if test coverage meets the minimum requirement of $(MIN_COVERAGE)%..."
	@go tool cover -func=coverage.out | grep 'total:' | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$1 >= $(MIN_COVERAGE)) {print "> Test coverage meets the minimum requirement: " $$1 "%"; exit_code=0} else {print "> [Error] Test coverage does not meet the minimum requirement: " $$1 "%. Required: " $(MIN_COVERAGE) "%"; exit_code=1}} END {exit(exit_code)}'

.PHONY: core-coverage
core-coverage:
	@echo "==> Running tests with coverage for the core layer (models and use cases)..."
	@go test -coverprofile=coverage.out $(shell go list ./... | grep -v /cmd | grep -v /config | grep -v /docs | grep -v /internal/database/ | grep -v /internal/webserver/ )
	@go tool cover -func=coverage.out
	@echo "> Generating HTML report..."
	@mkdir -p $(TEST_REPORTS_DIR)
	@go tool cover -html=coverage.out -o $(TEST_REPORTS_DIR)/coverage.html
	@echo "> Checking if test coverage meets the minimum requirement of $(MIN_CORE_COVERAGE)%..."
	@go tool cover -func=coverage.out | grep 'total:' | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$1 >= $(MIN_CORE_COVERAGE)) {print "> Test coverage meets the minimum requirement: " $$1 "%"; exit_code=0} else {print "> [Error] Test coverage does not meet the minimum requirement: " $$1 "%. Required: " $(MIN_CORE_COVERAGE) "%"; exit_code=1}} END {exit(exit_code)}'

clean:
	@echo "==> Cleaning up..."
	@rm -rf ./$(BUILD_DIR)
	@rm -rf ./$(TEST_REPORTS_DIR)
	@rm -f coverage.out
	@rm -f report-lint.html
	@go clean -testcache

build: clean
	@echo "==> Building the application..."
	@echo "> Target version: $(buildVersion)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/zg-data-guard/main.go
	@echo "version=$(buildVersion)" > $(BUILD_FILE)
	@echo "buildTime=$(shell date -u '+%Y-%m-%d %H:%M:%S')" >> $(BUILD_FILE)
	@echo "> Application built successfully at $(BUILD_DIR)/$(APP_NAME)!"

.PHONY: release
release: build
	@echo "==> Releasing the application..."
	@if [ -z "$${version}" ]; then \
		echo "> version is not set. Use: make release version=<version>. Release aborted."; \
		exit 1; \
	fi
	@echo "> Creating new tag: v$${version}"
	@git tag -a v$${version} -m "Release v$${version}"
	@echo "> Pushing new tag: v$${version}"
	@git push origin v$${version}
	@echo "> Release completed successfully!"

.PHONY: create_migration migrate_up migrate_down migrate_force run run-with-docs docs build test clean
