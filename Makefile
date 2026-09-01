#!/usr/bin/env make
# SPDX-FileCopyrightText: Copyright © 2025 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

export GOMOD_VER := 1.24.0
export GOTOOLCHAIN := go1.26.7+auto

export GOLANGCI_LINT_VER := latest
export GOFUMPT_VER := latest
export GORELEASER_VER := latest
export MISSPELL_VER := latest
export MODERNIZE_VER := latest
export VULN_VER := latest

# export GOLANGCI_LINT_VER := v2.12.2
# export GOFUMPT_VER := v0.11.0
# export GORELEASER_VER := v2.13.0
# export MISSPELL_VER := v0.8.0
# export MODERNIZE_VER := v0.48.0
# export VULN_VER := v1.2.0

export SHELL := /bin/bash
export NO_COLOR := 1
export TERM := dumb

TEST_TAGS :=$(strip $(TEST_TAGS),debug)

ifneq ($(wildcard go.tool.mod),)
TOOL_OPTS += -modfile=go.tool.mod
endif

export TOOL_OPTS

.DEFAULT_GOAL := all

.PHONY: all
all: ## make download gen build check test
all: download gen build check test

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
	rm -rf bin/*
	rm -rf dist
	rm -f *.bak
	rm -f coverage.*

# .PHONY: run
# run: ## go run
# 	go run .

.PHONY: mod
mod: ## go mod tidy
	go mod tidy -v -x $(GOMOD_VER)
	test -f go.tool.mod && go mod tidy -v -x $(TOOL_OPTS)

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
lint: ## golangci-lint run --fix .
	go tool $(TOOL_OPTS) golangci-lint run --fix ./...

.PHONY: vuln
vuln: ## govulncheck ./...
	-go tool $(TOOL_OPTS) govulncheck -show verbose ./...
	# We are pinned at go1.24.0, so we ignore issues for now

.PHONY: vet
vet: ## go vet ./...
	go vet .
	go vet github.com/rasa/compat/robustio
	go vet github.com/rasa/compat/volume
	go vet -unsafeptr=false ./...

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
test: ## go test .
	go test $(TEST_OPTS) -tags "$(TEST_TAGS)" $(RACE_OPT) -covermode=atomic -coverprofile=coverage.out -coverpkg=. ./...
	sed -i.bak "/compat\/cmd\//d; /compat\/golang\//d;" coverage.out
	rm -f *.bak
	go tool cover -html=coverage.out -o coverage.html

.PHONY: diff
diff: ## git diff
ifeq ($(OS),Windows_NT)
	git config --local core.filemode false
endif
	git --no-pager diff --exit-code
	@RES=$$(git status --porcelain --untracked-files=no) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

# Added by compat:

.PHONY: bld
bld: ## go build -trimpath -o ./bin/ ./...
	go build -trimpath -o ./bin/ ./...

.PHONY: bldtest
bldtest: ## make bld test
bldtest:         bld test

.PHONY: check
check: ## make fmt fumpt lint modernize spell vet restore diffx
check:         fmt fumpt lint modernize spell vet restore diffx

.PHONY: diffx
diffx: ## git diff -uw
	-git --no-pager diff -uw

.PHONY: download
download: ## go mod download
	go mod download
	test -f go.tool.mod && go mod download $(TOOL_OPTS)

.PHONY: fmt
fmt: ## go fmt ./...
	go fmt ./...

.PHONY: fumpt
fumpt: ## gofumpt -w .
	go tool $(TOOL_OPTS) gofumpt -w .

.PHONY: get
get: ## get gofumpt, golangci-lint, goreleaser, govulncheck, misspell, modernize
	export GOFLAGS="$(GOFLAGS) $(TOOL_OPTS)" ;\
	go get github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VER) ;\
	go get github.com/goreleaser/goreleaser/v2@$(GORELEASER_VER) ;\
	go get github.com/golangci/misspell/cmd/misspell@$(MISSPELL_VER) ;\
	go get golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@$(MODERNIZE_VER) ;\
	go get golang.org/x/vuln/cmd/govulncheck@$(VULN_VER) ;\
	go get mvdan.cc/gofumpt@$(GOFUMPT_VER)
	make mod

.PHONY: install
install: ## install gofumpt, golangci-lint, goreleaser, govulncheck, misspell, modernize
	export GOFLAGS="$(GOFLAGS) $(TOOL_OPTS)" ;\
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VER) ;\
	go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VER) ;\
	go install github.com/golangci/misspell/cmd/misspell@$(MISSPELL_VER) ;\
	go install golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@$(MODERNIZE_VER) ;\
	go install golang.org/x/vuln/cmd/govulncheck@$(VULN_VER) ;\
	go install mvdan.cc/gofumpt@$(GOFUMPT_VER)
	make mod

.PHONY: modernize
modernize: ## modernize ./...
	go tool $(TOOL_OPTS) modernize -fix ./...

.PHONY: nuke
nuke: ## make clean and go clean -i -cache -testcache -modcache -fuzzcache
nuke: clean
	rm -rf "$(shell go env GOCACHE)/../golangci-lint"
	go clean -i -cache -testcache -modcache -fuzzcache

.PHONY: release
release: ## Cut a new release via goreleaser
release:
	# go tool $(TOOL_OPTS) goreleaser release
	go tool $(TOOL_OPTS) goreleaser release

.PHONY: restore
restore: ##	git restore format.go walk.go walk_test.go golang/*.go robustio/robustio*.go
	git restore \
		format.go \
		walk.go \
		walk_test.go \
		golang/golang_*.go \
		golang/syscall/windows/*.go \
		golang/syscall/windows/registry/*.go \
		golang/syscall/windows/sysdll/*.go \
		golang/user/*.go \
		robustio/robustio*.go

.PHONY: update
update: ## go get -u
	go get -u
	test -f go.tool.mod && go get -u $(TOOL_OPTS)
	make mod

# aliases

.PHONY: tidy
tidy: mod

.PHONY: gofumpt
gofumpt: fumpt
