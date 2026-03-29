.PHONY: help run build test tidy fmt

help:
	@echo "Available targets:"
	@echo "  make run    - Run the game"
	@echo "  make build  - Build the gong binary"
	@echo "  make test   - Run Go tests"
	@echo "  make tidy   - Tidy go modules"
	@echo "  make fmt    - Format Go source files"

run:
	go run .

build:
	go build ./...

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w main.go game/*.go
