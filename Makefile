.PHONY: build test lint clean

build:
	go build ./...

test:
	go test ./...

lint:
	go vet ./...
	command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || true

clean:
	go clean ./...
