package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/grpc/protos"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDirectoryServiceServer
}

func (s *server) SendDirectoryStructure(ctx context.Context, in *pb.DirectoryRequest) (*pb.DirectoryResponse, error) {
	root := &pb.Directory{
		Name: "mock-root",
		Path: "/mock/path",
		Files: []*pb.File{
			{
				Name:         "hello.txt",
				OriginalPath: "/orig/hello.txt",
				NewPath:      "/mock/path/hello.txt",
				Tags: []*pb.Tag{
					{Name: "greeting"},
				},
				Metadata: []*pb.MetadataEntry{
					{Key: "size", Value: "1234"},
				},
			},
		},
		Directories: []*pb.Directory{
			{
				Name:        "subdir",
				Path:        "/mock/path/subdir",
				Files:       []*pb.File{},
				Directories: []*pb.Directory{},
			},
		},
	}

	return &pb.DirectoryResponse{Root: root}, nil
}

func main() {
	fmt.Println("running go server")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDirectoryServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
