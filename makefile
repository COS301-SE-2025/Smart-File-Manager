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
	pytest -v -s --color=yes --tb=short python/testing/

proto_gen:
	python -m grpc_tools.protoc \
  -Ipython/src/protos \
  --python_out=python/src \
  --pyi_out=python/src \
  --grpc_python_out=python/src \
  python/src/protos/helloworld.proto

python_server:
	python3 python/src/greeter_server.py

python_client:
	python3 python/src/greeter_client.py
