#!/usr/bin/env make
# SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

SHELL := /bin/bash
export NO_COLOR := 1
export TERM := dumb

TEST_TAGS :=$(strip $(TEST_TAGS),debug)

ifneq ($(wildcard go.tool.mod),)
TOOL_OPTS += -modfile=go.tool.mod
endif

export TOOL_OPTS

.DEFAULT_GOAL := all

.PHONY: all
all: ## make download gen build spell lint fix test
all: download gen build spell lint fix test vet

.PHONY: precommit
precommit: ## make all vuln
precommit: all vuln

.PHONY: ci
ci: ## make precommit diff
ci: precommit diff

.PHONY: help
help:
	@awk -F ':.*##[ \t]*' '/^[^#: \t]+:.*##/ {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: clean
clean: ## remove files created during build pipeline
	rm -rf dist
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: run
run: ## go run
	go run .

.PHONY: mod
mod: ## go mod tidy
	go mod tidy
	test -f go.tool.mod && go mod tidy $(TOOL_OPTS)

.PHONY: gen
gen: ## go generate ./...
	go generate ./...

.PHONY: build
build: ## goreleaser build --clean --single-target --snapshot
	-go tool $(TOOL_OPTS) goreleaser --version
	go tool $(TOOL_OPTS) goreleaser build --clean --single-target --snapshot

.PHONY: spell
spell: ## misspell -error -locale=US -w **.md
	go tool $(TOOL_OPTS) misspell -error -locale=US -w **.md

.PHONY: lint
lint: ## golangci-lint run --fix
	go tool $(TOOL_OPTS) golangci-lint run --fix

.PHONY: fix
fix: ## gofumpt
	go tool $(TOOL_OPTS) gofumpt -w .
	git restore walk.go walk_test.go golang/golang_*.go

.PHONY: vuln
vuln: ## govulncheck
	go tool $(TOOL_OPTS) govulncheck ./...

.PHONY: modernize
modernize: ## modernize
	go tool $(TOOL_OPTS) modernize -fix ./...

.PHONY: vet
vet: ## vet
	go vet ./...

# Added by compat:

.PHONY: download
download: ## go mod download
	go mod download
	test -f go.tool.mod && go mod download $(TOOL_OPTS)
	# make mod

.PHONY: get
get: ## go get -u
	go get -u
	test -f go.tool.mod && go get -u $(TOOL_OPTS)
	make mod

.PHONY: tools
tools: ## freshen tools (misspell, golangci-lint, goreleaser, govulncheck, gofumpt)
	export GOFLAGS="$(GOFLAGS) $(TOOL_OPTS)" ;\
	go get github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest ;\
	go get github.com/goreleaser/goreleaser/v2@latest ;\
	go get golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize ;\
	go get github.com/client9/misspell/cmd/misspell@latest ;\
	go get golang.org/x/vuln/cmd/govulncheck@latest ;\
	go get mvdan.cc/gofumpt@latest
	make mod

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
HAS_CC := $(shell command -v $(CC) >/dev/null)
ifeq ($(HAS_CC),)
RACE_OPT =
endif

.PHONY: test
test: ## go test
	go test $(TEST_OPTS) -tags "$(TEST_TAGS)" $(RACE_OPT) -covermode=atomic -coverprofile=coverage.out -coverpkg=. .
	sed -i.bak "/compat\/cmd\//d; /compat\/golang\//d;" coverage.out
	rm -f *.bak
	go tool cover -html=coverage.out -o coverage.html

.PHONY: diff
diff: ## git diff
ifeq ($(OS),Windows_NT)
	git config --local core.filemode false
endif
	git diff --exit-code
	@RES=$$(git status --porcelain --untracked-files=no) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi
