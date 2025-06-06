# Use Golang base image
FROM golang:1.24.3-alpine

RUN apk add --no-cache protobuf protoc \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set working directory inside container
WORKDIR /app

# Copy over files
COPY golang/ ./golang/
COPY protos/     ./protos/

# Copy over folder for tests
COPY testRootFolder ./testRootFolder

# Use app as PATH
ENV GOPATH=/app

# Required since we installed protoc
ENV PATH="/go/bin:${PATH}"

# Install golang dependancies (None yet)

# Generate gRPC code
RUN protoc \
	--proto_path=protos \
	--go_out=golang \
	--go-grpc_out=golang \
	protos/message_structure.proto

WORKDIR /app/golang

# Expose the API server port
EXPOSE 51000

# Start the server for API requests
CMD ["go", "run", "main.go"]

# Old command that used to run tests directly in docker (Comment this out for now as the workflow will do it)
# CMD ["go", "test", "-v", "./filesystem"]