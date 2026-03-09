.PHONY: build run test test-v test-race cover cover-html lint vet fmt tidy clean docker-up docker-down docker-build

APP_NAME := evv-logger
BUILD_DIR := bin
COVER_FILE := coverage.out

## build: compile the binary
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/http

## run: start the server locally
run:
	go run ./cmd/http

## test: run all tests
test:
	go test ./...

## test-v: run all tests with verbose output
test-v:
	go test -v ./...

## test-race: run all tests with race detector
test-race:
	go test -race ./...

## cover: run tests and print coverage summary
cover:
	go test -coverprofile=$(COVER_FILE) ./...
	go tool cover -func=$(COVER_FILE)

## cover-html: open coverage report in browser
cover-html: cover
	go tool cover -html=$(COVER_FILE)

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## vet: run go vet
vet:
	go vet ./...

## fmt: format all Go files
fmt:
	gofmt -w .

## tidy: tidy and verify module dependencies
tidy:
	go mod tidy
	go mod verify

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR) $(COVER_FILE)

## docker-build: build the Docker image
docker-build:
	docker compose -f deployment/docker/docker-compose.yaml build

## docker-up: start containers
docker-up:
	docker compose -f deployment/docker/docker-compose.yaml up --build -d

## docker-down: stop and remove containers
docker-down:
	docker compose -f deployment/docker/docker-compose.yaml down

## check: run fmt, vet, lint, and tests (pre-commit check)
check: fmt vet lint test-race

## help: print this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
