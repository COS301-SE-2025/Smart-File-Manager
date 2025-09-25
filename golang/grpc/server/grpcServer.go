//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
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

func loadEnvFile(path string) (map[string]string, error) {
	vars := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// remove "export " if present
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars, nil
}
func FindProjectRoot(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // reached filesystem root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("%s not found", filename)
}

func main() {
	fmt.Println("Starting go server...")
	// Read port from environment
	path, err := FindProjectRoot("server.env")
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	env, err := loadEnvFile(path)
	if err != nil {
		log.Fatalf("failed to read env file: %v", err)
	}

	port := env["PYTHON_SERVER"]

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	pb.RegisterDirectoryServiceServer(s, &server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
