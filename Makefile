.PHONY: help run build test tidy fmt web serve-web

help:
	@echo "Available targets:"
	@echo "  make run    - Run the game"
	@echo "  make build  - Build the gong binary"
	@echo "  make test   - Run Go tests"
	@echo "  make tidy   - Tidy go modules"
	@echo "  make fmt    - Format Go source files"
	@echo "  make web    - Build the WebAssembly browser game"
	@echo "  make serve-web - Build and serve the browser game locally"

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

web:
	rm -rf dist
	mkdir -p dist
	GOOS=js GOARCH=wasm go build -o dist/gong.wasm .
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" dist/wasm_exec.js
	cp web/index.html web/styles.css web/app.js dist/

serve-web: web
	python3 -m http.server 8080 --directory dist
