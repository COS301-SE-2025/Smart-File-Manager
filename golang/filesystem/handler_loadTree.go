package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	//grpc imports
	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
	"google.golang.org/grpc"
)

// struct for json return type for 200 reqs
type DirectoryTreeJson struct {
	Name     string     `json:"name"`
	IsFolder bool       `json:"isFolder"`
	Children []FileNode `json:"children"`
}

// file or folder
type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path,omitempty"`
	IsFolder bool       `json:"isFolder"`
	Tags     []string   `json:"tags,omitempty"`
	Metadata *Metadata  `json:"metadata,omitempty"`
	Children []FileNode `json:"children,omitempty"`
}

type Metadata struct {
	Size         string `json:"size"`
	DateCreated  string `json:"dateCreated"`
	Owner        string `json:"owner"`
	LastModified string `json:"lastModified"`
}

func grpcFunc(c *Folder) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create new grpc go client")
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DirectoryRequest{
		Root: convertFolderToProto(*c),
	}

	resp, err := client.SendDirectoryStructure(ctx, req)
	if err != nil {
		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	}

	// fmt.Printf("Server returned root directory: name=%q, path=%q/n", resp.Root.GetName(), resp.Root.GetPath())
	// printDirectory(resp.Root, 0)

	printDirectoryWithMetadata(resp.Root, 0)
}

func convertFolderToProto(f Folder) *pb.Directory {
	protoDir := &pb.Directory{
		Name: f.Name, // make sure ItemName is exported
		Path: f.Path,
	}

	for _, file := range f.Files {
		protoDir.Files = append(protoDir.Files, &pb.File{
			Name:         file.Name,
			OriginalPath: file.Path,
		})
		// case *fs.Folder:
		// 	protoDir.Directories = append(protoDir.Directories, convertFolderToProto(v))

	}
	for _, subFolder := range f.Subfolders {
		protoDir.Directories = append(protoDir.Directories, convertFolderToProto(*subFolder))

	}
	return protoDir
}

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

// actual api endpoint function
func loadTreeDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		if c.Name == name {
			// build the nested []FileNode
			children := createDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				Children: children,
			}

			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			grpcFunc(c)
			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

// helper recursive function
func createDirectoryJSONStructure(folder *Folder) []FileNode {
	var nodes []FileNode

	for _, file := range folder.Files {

		md := Metadata{}
		tags := file.Tags
		if len(tags) == 0 {
			tags = []string{"none"}
		}

		nodes = append(nodes, FileNode{
			Name:     file.Name,
			Path:     file.Path,
			IsFolder: false,
			Tags:     tags,
			Metadata: &md,
		})
	}

	for _, sub := range folder.Subfolders {
		// recurse first
		childNodes := createDirectoryJSONStructure(sub)

		nodes = append(nodes, FileNode{
			Name:     sub.Name,
			Path:     sub.Path,
			IsFolder: true,
			Tags:     sub.Tags,
			Metadata: &Metadata{},
			Children: childNodes,
		})
	}

	return nodes
}
