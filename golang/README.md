go get -u github.com/golang/protobuf/
go get -u google.golang.org/grpc

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest