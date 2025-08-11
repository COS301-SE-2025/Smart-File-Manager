package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	RootPath string     `json:"rootPath"`
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
	Locked   bool       `json:"locked"`
	NewPath  string     `json:"newPath,omitempty"` // for moving files
}

type Metadata struct {
	Size         string `json:"size"`
	DateCreated  string `json:"dateCreated"`
	MimeType     string `json:"mimeType"`
	LastModified string `json:"lastModified"`
}

func grpcFunc(c *Folder, requestType string) error {
	// fmt.Println("+++++pretty print start+++++")
	// printFolderDetails(c, 0)
	// fmt.Println("+++++pretty print end+++++")
	if requestType != "METADATA" && requestType != "CLUSTERING" {
		return fmt.Errorf("invalid requestType: %s", requestType)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create new grpc go client")
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	req := &pb.DirectoryRequest{
		Root:        convertFolderToProto(*c),
		RequestType: requestType,
	}

	resp, err := client.SendDirectoryStructure(ctx, req)
	if err != nil {
		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	}

	fmt.Printf("Server returned root directory: name=%q, path=%q/n", resp.Root.GetName(), resp.Root.GetPath())

	mergeProtoToFolder(resp.Root, c)

	return nil
}

func mergeProtoToFolder(dir *pb.Directory, existing *Folder) {
	if dir == nil || existing == nil {

		fmt.Println("mergeProtoToFolder called with nil")
		return
	}

	existing.Files = nil
	existing.Subfolders = nil

	for _, file := range dir.Files {
		existing.Files = append(existing.Files, &File{
			Name:     file.Name,
			Path:     file.OriginalPath,
			Tags:     tagsToStrings(file.Tags),
			Metadata: metadataConverter(file.Metadata),
			Locked:   file.IsLocked,
			NewPath:  file.NewPath,
		})
	}

	for _, sub := range dir.Directories {
		child := &Folder{
			Name:   sub.Name,
			Path:   sub.Path,
			Locked: sub.IsLocked,
		}
		mergeProtoToFolderHelper(sub, child)
		existing.Subfolders = append(existing.Subfolders, child)
	}
}

func mergeProtoToFolderHelper(dir *pb.Directory, existing *Folder) {
	if dir == nil || existing == nil {

		fmt.Println("mergeProtoToFolderHelper called with nil")
		return
	}
	existing.Name = dir.Name
	existing.Path = dir.Path

	existing.Files = nil
	existing.Subfolders = nil

	for _, file := range dir.Files {
		existing.Files = append(existing.Files, &File{
			Name:     file.Name,
			Path:     file.OriginalPath,
			Tags:     tagsToStrings(file.Tags),
			Metadata: metadataConverter(file.Metadata),
			Locked:   file.IsLocked,
			NewPath:  file.NewPath,
		})
	}

	for _, sub := range dir.Directories {
		child := &Folder{
			Name:   sub.Name,
			Path:   sub.Path,
			Locked: sub.IsLocked,
		}
		mergeProtoToFolderHelper(sub, child)
		existing.Subfolders = append(existing.Subfolders, child)
	}
}

func tagsToStrings(tags []*pb.Tag) []string {
	var tagStrings []string

	for _, tag := range tags {
		fmt.Println("TAG FOUND: " + tag.Name)
		tagStrings = append(tagStrings, tag.Name)
	}
	return tagStrings
}

// converts string version of tags to *pb.Tag array
func stringsToTags(stringTags []string) []*pb.Tag {
	var tags []*pb.Tag

	for _, tag := range stringTags {
		fmt.Println("TAG FOUND: " + tag)
		curr := &pb.Tag{
			Name: tag,
		}
		tags = append(tags, curr)
	}
	return tags
}

// converts the compositite structure to the correct structure grpc uses
func convertFolderToProto(f Folder) *pb.Directory {
	protoDir := &pb.Directory{
		Name: f.Name, // make sure ItemName is exported
		Path: f.Path,
	}

	for _, file := range f.Files {
		protoDir.Files = append(protoDir.Files, &pb.File{
			Name:         file.Name,
			OriginalPath: file.Path,
			Tags:         stringsToTags(file.Tags),
			IsLocked:     file.Locked,
			NewPath:      file.NewPath,
		})
		// case *fs.Folder:
		// 	protoDir.Directories = append(protoDir.Directories, convertFolderToProto(v))

	}
	for _, subFolder := range f.Subfolders {
		protoDir.Directories = append(protoDir.Directories, convertFolderToProto(*subFolder))

	}
	return protoDir
}

// pb.metadataentry[] to the filesystem metadataEntry[]
func metadataConverter(metaDataArr []*pb.MetadataEntry) []*MetadataEntry {
	// preallocate a slice of the correct length
	res := make([]*MetadataEntry, len(metaDataArr))

	// copy each protobuf entry into local type
	for i, entry := range metaDataArr {
		res[i] = &MetadataEntry{
			Key:   entry.Key,
			Value: entry.Value,
		}
	}

	return res
}

// convert filesystem metadataEntry[] to Metadata struct for json response
func extractMetadata(metaDataArr []*MetadataEntry) *Metadata {
	// fmt.Println("extractMetadata called. metaDataArr len: " + strconv.Itoa(len(metaDataArr)))

	md := &Metadata{}

	fieldMap := map[string]*string{
		"size_bytes": &md.Size,
		"created":    &md.DateCreated,
		"mime_type":  &md.MimeType,
		"modified":   &md.LastModified,
	}

	for _, entry := range metaDataArr {
		// fmt.Println(entry)
		if ptr, ok := fieldMap[entry.Key]; ok {
			*ptr = entry.Value
		}
	}
	return md
}

func printDirectoryWithMetadata(dir *pb.Directory, num int) {

	space := strings.Repeat("  ", num)
	fmt.Println(space + "Directory: " + dir.Name)
	fmt.Println(space + "Path: " + dir.Path)

	for _, file := range dir.Files {
		fmt.Println(space + "File name: " + file.Name)
		fmt.Println(space + "File original path: " + file.OriginalPath)
		fmt.Println(space + "File new path: " + file.NewPath)

		fmt.Println("----")
	}

	// for _, dir := range dir.Directories {
	// 	newNum := num + 1
	// 	printDirectoryWithMetadata(dir, newNum)
	// }

}

// endpoint called using no grpc:
func loadTreeDataHandlerGoOnly(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GOVERSION OF loadTree called")
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		if c.Name == name {

			children := GoSidecreateDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				RootPath: c.Path,
				Children: children,
			}
			// PrettyPrintFolder(c, "")

			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			LoadTypes(c, name) // load types into the global objectMap
			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)

}

// helper recursive function
func GoSidecreateDirectoryJSONStructure(folder *Folder) []FileNode {
	var nodes []FileNode

	for _, file := range folder.Files {

		fi, err := os.Stat(file.Path)

		md := &Metadata{}

		if err != nil {
			fmt.Println(err)
			md = nil
		} else {
			layout := "2006-01-02 15:04"
			md.Size = strconv.FormatInt(fi.Size(), 10)

			md.DateCreated = fi.ModTime().Format(layout)
			md.LastModified = fi.ModTime().Format(layout)

			md.MimeType = ""
			lastDotIndex := strings.LastIndex(file.Name, ".")
			if lastDotIndex != -1 {
				// Slice the string from the last period's index to the end
				md.MimeType = file.Name[lastDotIndex:]
			} else {
				md.MimeType = "mystery"
			}

			mdEntries := []*MetadataEntry{
				{Key: "Size", Value: strconv.FormatInt(fi.Size(), 10)},
				{Key: "DateCreated", Value: fi.ModTime().Format(layout)},
				{Key: "LastModified", Value: fi.ModTime().Format(layout)},
			}
			file.Metadata = mdEntries
			// file.Tags =

		}

		tags := file.Tags

		nodes = append(nodes, FileNode{
			Name:     file.Name,
			Path:     file.Path,
			IsFolder: false,
			Tags:     tags,
			Metadata: md,
			Locked:   file.Locked,
		})
	}

	for _, sub := range folder.Subfolders {
		// recurse first
		childNodes := GoSidecreateDirectoryJSONStructure(sub)

		nodes = append(nodes, FileNode{
			Name:     sub.Name,
			Path:     sub.Path,
			IsFolder: true,
			Tags:     sub.Tags,
			Metadata: &Metadata{},
			Children: childNodes,
			Locked:   sub.Locked,
		})
	}

	return nodes
}

func createDirectoryJSONStructure(folder *Folder) []FileNode {
	var nodes []FileNode

	for _, file := range folder.Files {

		md := extractMetadata(file.Metadata)
		tags := file.Tags
		// if len(tags) == 0 {
		// 	tags = []string{"none"}
		// }

		nodes = append(nodes, FileNode{
			Name:     file.Name,
			Path:     file.Path,
			IsFolder: false,
			Tags:     tags,
			Metadata: md,
			Locked:   file.Locked,
			NewPath:  file.NewPath, // include NewPath for moving files
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
			Locked:   sub.Locked,
			NewPath:  sub.NewPath,
		})
	}

	return nodes
}
