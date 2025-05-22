# .PHONY: test lint fmt build
# fmt:
# 	go fmt ./...
# 	ruff format python/
# lint:
# 	golangci-lint run ./...
# 	ruff python/
# test:
# 	go test ./...
# 	pytest -q
# build:
# 	go build -o bin/sfm ./go/cmd/sfm

# Otherwise make thinks these are files and not commands
.PHONY: python python_test

go:
	go run ./golang/main.go

python:
	python3 python/src/main.py

python_test:
	pytest -v --color=yes --tb=short python/testing/