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
.PHONY: python python_test go_test

go:
	go run ./golang/main.go

go_test:
	cd golang/filesystem && go test -v

python:
	python3 python/src/main.py

python_test:
	pytest -v -s --color=yes --tb=short python/testing/

proto_gen:
	python3 -m grpc_tools.protoc \
		-Iprotos \
		--python_out=python/src \
		--pyi_out=python/src \
		--grpc_python_out=python/src \
		protos/message_structure.proto

	protoc \
	  --proto_path=protos \
	  --go_out=golang \
	  --go-grpc_out=golang \
	  protos/message_structure.proto


python_server:
	python3 python/src/request_handler.py

python_client:
	python3 python/src/greeter_client.py
