build:
	@go build -o bin/559 cmd/559/main.go

buildw:
	@go build -o bin/559.exe cmd/559/main.go

run:
	@go run cmd/559/main.go


test:
	go test ./...