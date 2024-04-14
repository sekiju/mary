.PHONY: build test

ifeq ($(OS),Windows_NT)
    EXE_SUFFIX := .exe
else
    EXE_SUFFIX :=
endif

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	go build -ldflags '-s -w' -o build/mary$(EXE_SUFFIX) cmd/cli/main.go

run:
	go run cmd/cli/main.go

test:
	go test ./internal/connectors/...

test-clean:
	go clean -testcache
