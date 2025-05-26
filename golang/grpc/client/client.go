package main

import (
	"context"
	"fmt"
	"log"
	"strings"
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

	fmt.Printf("Server returned root directory: name=%q, path=%q\n", resp.Root.GetName(), resp.Root.GetPath())
	printDirectory(resp.Root, 0)
}

func printDirectory(dir *pb.Directory, num int) {
	fmt.Println("Directory: " + dir.Name)
	fmt.Println("Path: " + dir.Path)
	space := strings.Repeat("  ", num)
	for _, file := range dir.Files {
		fmt.Println(space + "File name: " + file.Name)
		fmt.Println(space + "File origional path: " + file.OriginalPath)
		fmt.Println("----")
	}
	for _, dir := range dir.Directories {
		newNum := num + 1
		printDirectory(dir, newNum)
	}

}
