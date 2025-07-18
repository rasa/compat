SHELL := /bin/bash

.DEFAULT_GOAL := all

.PHONY: all
all: ## build pipeline
all: mod gen build spell lint test

.PHONY: precommit
precommit: ## validate the branch before commit
precommit: all vuln

.PHONY: ci
ci: ## CI build pipeline
ci: precommit diff

.PHONY: help
help:
	@awk -F ':.*##[ \t]*' '/^[^#: \t]+:.*##/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: clean
clean: ## remove files created during build pipeline
	rm -rf dist
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: mod
mod: ## go mod tidy
	go mod tidy

.PHONY: gen
gen: ## go generate
	go generate ./...

.PHONY: build
build: ## goreleaser build
	go tool goreleaser build --clean --single-target --snapshot

.PHONY: spell
spell: ## misspell
	go tool misspell -error -locale=US -w **.md

.PHONY: lint
lint: ## golangci-lint
	go tool golangci-lint run --fix

.PHONY: vuln
vuln: ## govulncheck
	go tool govulncheck ./...

RACE_OPT := -race

# go: -race requires cgo
ifneq ($(strip $(CGO_ENABLED)),1)
RACE_OPT =
endif

GO_VERSION := $(shell go version)
# go: -race is not supported on windows/arm64
ifeq ($(findstring windows/arm64,$(GO_VERSION)),windows/arm64)
RACE_OPT =
endif

# cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
CC := $(shell go env CC)
HAS_CC := $(shell command -v $(CC) >/dev/null || echo no)
ifeq ($(HAS_CC),no)
RACE_OPT =
endif

.PHONY: test
test: ## go test
	go test $(TEST_OPTS) -tags debug $(RACE_OPT) -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: diff
diff: ## git diff
	git diff --exit-code
	@RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi
