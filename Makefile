build:
	@go build -o bin/559 cmd/559/main.go

test:
	go test ./...