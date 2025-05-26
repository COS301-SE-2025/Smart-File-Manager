package main

import (
	"fmt"
	"log"
	"os"

	fs "github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	rootPath := "../testRootFolder"

	// Check that the test directory exists before continuing
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist. Please create it before running the test.", rootPath)
	}

	// Convert root directory to composite
	root := fs.ConvertToComposite("001", "TestRoot", rootPath)

	// Simple validation
	if len(root.ContainedItems) == 0 {
		log.Fatalf("Expected at least one item in the root folder, but got none.")
	}

	// 	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 	if err != nil {
	// 		log.Fatalf("could not create new grpc go client")
	// 	}
	// 	defer conn.Close()

	// 	client := pb.NewDirectoryServiceClient(conn)

	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()

	// 	routePath := filepath.Clean(
	// 		`/mnt/c/Users/jackb`,
	// 	)

	// 	req := &pb.DirectoryRequest{
	// 		Root: &pb.Directory{
	// 			Name: "jackb",
	// 			Path: routePath,
	// 		},
	// 	}

	// 	resp, err := client.SendDirectoryStructure(ctx, req)
	// 	if err != nil {
	// 		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	// 	}
	// 	fmt.Println("++++++++++++++++++++++++++")
	// 	fmt.Println("route path: " + routePath)
	// 	fmt.Println("end")

	// 	// fmt.Printf("Server returned root directory: name=%q, path=%q/n", resp.Root.GetName(), resp.Root.GetPath())
	// 	printDirectory(resp.Root, 0)
	// }

	// func printDirectory(dir *pb.Directory, num int) {

	// 	fmt.Println("Directory: " + dir.Name)
	// 	fmt.Println("Path: " + dir.Path)

	// 	space := strings.Repeat("  ", num)
	// 	for _, file := range dir.Files {
	// 		fmt.Println(space + "File name: " + file.Name)
	// 		fmt.Println(space + "File origional path: " + file.OriginalPath)
	// 		fmt.Println("----")
	// 	}

	// 	for _, dir := range dir.Directories {
	// 		newNum := num + 1
	// 		printDirectory(dir, newNum)
	// 	}

}
