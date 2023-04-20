.DEFAULT_GOAL := help
SHELL := /bin/bash

APP_NAME := envctl
BUILD_DIR := ./dist
DST_DIR := /usr/local/bin

CURRENT_SHA = $(shell git rev-parse --short HEAD)
CURRENT_TS = $(shell date -I seconds)

##@ App

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint: ## Lint the app
	actionlint
	golangci-lint run

.PHONY: format
format: ## Format the app
	gofmt -w .
	pin-github-action .github/workflows/*.yml

.PHONY: test
test: ## Test the app
	go mod tidy
	go test --coverprofile=coverage.out ./...

.PHONY: coverage
coverage: ## Show code coverage
	@make test
	gocovsh --profile coverage.out

.PHONY: build
build: ## Build the app
	go mod tidy
	go build -ldflags='-X main.version=dev -X main.commit=$(CURRENT_SHA) -X main.date=$(CURRENT_TS)' -o ${BUILD_DIR}/${APP_NAME} .

.PHONY: install
install: build ## Install the app
	install -d ${DST_DIR}
	install -m755 ${BUILD_DIR}/${APP_NAME} ${DST_DIR}/

.PHONY: version
version: ## Calculate the next release version
	./bin/version.sh

.PHONY: release
release: ## Create a new release tag
	./bin/release.sh


##@ Other

.PHONY: setup
setup: ## Bootstrap for local development
	./bin/setup.sh

# Via https://www.thapaliya.com/en/writings/well-documented-makefiles/
# Note: The `##@` comments determine grouping
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
