.PHONY: build test

build:
	go build -ldflags '-s -w' -o build/mary.exe cmd/cli/main.go

run:
	go run cmd/cli/main.go

test:
	go test ./internal/connectors/...

test-clean:
	go clean -testcache
