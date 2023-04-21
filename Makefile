.DEFAULT_GOAL := help
SHELL := /bin/bash

ARTIFACT_PATH = $(shell jq -r '.[0].path' dist/artifacts.json)
INSTALL_DIR = /usr/local/bin


##@ App

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint: ## Lint the app
	stylist check

.PHONY: format
format: ## Format the app
	stylist fix

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
	goreleaser build --clean --snapshot --single-target

.PHONY: install
install: build ## Install the app
	install -d ${INSTALL_DIR}
	install -m755 "${ARTIFACT_PATH}" ${INSTALL_DIR}/
	du -h "${ARTIFACT_PATH}"

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
