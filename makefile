.PHONY: test lint fmt build
fmt:
	go fmt ./...
	ruff format python/
lint:
	golangci-lint run ./...
	ruff python/
test:
	go test ./...
	pytest -q
build:
	go build -o bin/sfm ./go/cmd/sfm
