//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"

	fs "github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// rootPath := "/mnt/c/Users/jackb/OneDrive - University of Pretoria/Documents/TUKS/year 3/COS301/capstone/Smart-File-Manager/python/testing"

	path, _ := os.Getwd()
	fmt.Println("THE PATH: " + path)
	path = filepath.Dir(path)
	path = filepath.Join(path, "python/testing")
	fmt.Println("THE PATH: " + path)
	rootPath := path

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist. Please create it before running the test.", rootPath)
	}

	root := fs.ConvertToComposite("001", "TestRoot", rootPath)

	if len(root.ContainedItems) == 0 {
		log.Fatalf("Expected at least one item in the root folder, but got none.")
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create new grpc go client")
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DirectoryRequest{
		Root: convertFolderToProto(root),
	}

	resp, err := client.SendDirectoryStructure(ctx, req)
	if err != nil {
		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	}

	// fmt.Printf("Server returned root directory: name=%q, path=%q/n", resp.Root.GetName(), resp.Root.GetPath())
	// printDirectory(resp.Root, 0)

	printDirectoryWithMetadata(resp.Root, 0)

}

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

// }

func printDirectoryWithMetadata(dir *pb.Directory, num int) {

	space := strings.Repeat("  ", num)
	fmt.Println(space + "Directory: " + dir.Name)
	fmt.Println(space + "Path: " + dir.Path)

	for _, file := range dir.Files {
		fmt.Println(space + "File name: " + file.Name)
		fmt.Println(space + "File original path: " + file.OriginalPath)
		fmt.Println(space + "=====METADATA=========")
		metaData := file.Metadata

		for _, singleData := range metaData {
			fmt.Println(space + singleData.Key + "  :  " + singleData.Value)
		}
		fmt.Println("----")
	}

	for _, dir := range dir.Directories {
		newNum := num + 1
		printDirectoryWithMetadata(dir, newNum)
	}

}

// convertFolderToProto recursively builds a *pb.Directory from our *fs.Folder
func convertFolderToProto(f *fs.Folder) *pb.Directory {
	protoDir := &pb.Directory{
		Name: f.ItemName, // make sure ItemName is exported
		Path: f.GetPath(),
	}

	for _, child := range f.ContainedItems {
		switch v := child.(type) {
		case *fs.File:
			protoDir.Files = append(protoDir.Files, &pb.File{
				Name:         v.ItemName,
				OriginalPath: v.GetPath(),
			})
		case *fs.Folder:
			protoDir.Directories = append(protoDir.Directories, convertFolderToProto(v))
		}
	}
	return protoDir
}
