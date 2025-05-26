package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create new grpc go client")
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DirectoryRequest{
		Root: &pb.Directory{
			Name: "client-root",
			Path: "/tmp",
		},
	}
	resp, err := client.SendDirectoryStructure(ctx, req)
	if err != nil {
		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	}

	// 5) Print out whatever the server sent back
	fmt.Printf("Server returned root directory: name=%q, path=%q\n",
		resp.Root.GetName(), resp.Root.GetPath())
}
