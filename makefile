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
	go run grpc/server/grpcServer.go

go_grpc_client:
	cd golang && \
	go run grpc/client/grpcClient.go

go_test:
	cd golang/filesystem && go test -tags=test -v

go_api:
	cd golang && \
	go run .

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

	    mkdir -p golang/client
	    protoc -I. \
	    --go_out=golang/client \
	    --go_opt=paths=source_relative \
	    --go_opt=Mprotos/message_structure.proto=github.com/COS301-SE-2025/Smart-File-Manager/golang/client \
	    --go-grpc_out=golang/client \
	    --go-grpc_opt=paths=source_relative \
	    --go-grpc_opt=Mprotos/message_structure.proto=github.com/COS301-SE-2025/Smart-File-Manager/golang/client \
	    protos/message_structure.proto

python_client:
	python3 python/src/greeter_client.py

python_master_temp:
	pytest -v -s --color=yes --tb=short python/testing/test_clustering_request.py

