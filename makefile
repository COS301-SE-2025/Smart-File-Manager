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
	@echo "Running filesystem tests..."
	cd golang/filesystem && go test -tags=test -v
	@echo "Running filesystem/tests..."
	cd golang/filesystem/tests && go test -tags=test -v

go_coverage:
	cd golang/filesystem && go test -coverprofile=coverage.out -covermode=atomic
	@echo "Coverage summary:"
	cd golang/filesystem && go tool cover -html=coverage.out
	@echo "To view HTML report, run: go tool cover -html=golang/filesystem/coverage.out"

go_api:
	cd golang && \
	go run .

python:
	python3 python/src/main.py

python_test:
	pytest -v -s --color=yes --tb=short python/testing/

python_test_pyinstrument:
	pyinstrument -r html -o profiling/profile_report.html -m pytest -v -s --color=yes --tb=short python/testing/

python_test_clustering_request_pyinstrument:
	pyinstrument --renderer html -o profile.html -m pytest -v -s --color=yes --tb=short python/testing/test_clustering_request.py

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

python_locked_temp:
	pytest -v -s --color=yes --tb=short python/testing/test_locked_request.py

python_fn_temp:
	pytest -v -s --color=yes --tb=short python/testing/test_folder_name_creator.py
