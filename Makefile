.DEFAULT_GOAL := help
SHELL := /bin/bash

BIN_NAME = $(shell jq -r '.project_name' dist/metadata.json)
CMP_NAME = ${BIN_NAME}.bash
MAN_NAME = ${BIN_NAME}.1.gz

BIN_BUILD_PATH = $(shell jq -r '.[0].path' dist/artifacts.json)
CMP_BUILD_PATH = build/completions/${CMP_NAME}
MAN_BUILD_PATH = build/manpages/${MAN_NAME}

BIN_INSTALL_DIR = /usr/local/bin
CMP_INSTALL_DIR = $(shell brew --prefix)/etc/bash_completion.d
MAN_INSTALL_DIR = /usr/local/share/man/man1

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
	du -h "${BIN_BUILD_PATH}"

.PHONY: install
install: build ## Install the app
	install -d ${BIN_INSTALL_DIR}
	install -m755 "${BIN_BUILD_PATH}" ${BIN_INSTALL_DIR}/
	install -d ${CMP_INSTALL_DIR}
	install -m755 "${CMP_BUILD_PATH}" ${CMP_INSTALL_DIR}/
	install -d ${MAN_INSTALL_DIR}
	install -m644 "${MAN_BUILD_PATH}" ${MAN_INSTALL_DIR}/

.PHONY: uninstall
uninstall: ## Uninstall the app
	rm -f ${BIN_INSTALL_DIR}/${BIN_NAME}
	rm -f ${CMP_INSTALL_DIR}/${CMP_NAME}
	rm -f ${MAN_INSTALL_DIR}/${MAN_NAME}

.PHONY: version
version: ## Calculate the next release version
	./bin/version.sh

.PHONY: release
release: ## Create a new release tag
	./bin/release.sh

.PHONY: goreleaser
goreleaser:
	goreleaser release --snapshot --clean


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
