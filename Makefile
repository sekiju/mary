build:
	go build -o bin/mary cmd/cli/main.go

buildw:
	go build -o bin/mary.exe cmd/cli/main.go

run:
	go run cmd/cli/main.go

test:
	go test ./internal/connectors/...

test-cover:
	go test ./internal/connectors/... -v -coverprofile cover.out

update-all:
	go get -u ./... && go mod tidy