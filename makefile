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

go_proto_gen:
	mkdir -p golang/client
	protoc -I. \
	--go_out=golang/client \
	--go_opt=paths=source_relative \
	--go_opt=Mprotos/message_structure.proto=github.com/COS301-SE-2025/Smart-File-Manager/golang/client \
	--go-grpc_out=golang/client \
	--go-grpc_opt=paths=source_relative \
	--go-grpc_opt=Mprotos/message_structure.proto=github.com/COS301-SE-2025/Smart-File-Manager/golang/client \
	protos/message_structure.proto

go_grpc_server:
	cd golang && \
	go run grpc/server/server.go

go_grpc_client:
	cd golang && \
	go run grpc/client/client.go


python:
	python3 python/src/main.py

python_test:
	pytest -v -s --color=yes --tb=short python/testing/

proto_gen:
	python -m grpc_tools.protoc \
		-Iprotos \
		--python_out=python/src \
		--pyi_out=python/src \
		--grpc_python_out=python/src \
		protos/message_structure.proto

#	protoc \
		--go_out=golang \
		--go-grpc_out=golang \
		--proto_path=protos \
		protos/helloworld.proto


python_server:
	python3 python/src/request_handler.py

python_client:
	python3 python/src/greeter_client.py
