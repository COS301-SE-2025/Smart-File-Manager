package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	Name     string        `json:"name"`
	Path     string        `json:"path,omitempty"`
	IsFolder bool          `json:"isFolder"`
	Tags     []string      `json:"tags,omitempty"`
	Metadata *Metadata     `json:"metadata,omitempty"`
	Children []FileNode    `json:"children,omitempty"`
	Keywords []*pb.Keyword `json:"keywords,omitempty"`
	Locked   bool          `json:"locked"`
	NewPath  string        `json:"newPath,omitempty"` // for moving files
}

type Metadata struct {
	Size         string `json:"size"`
	DateCreated  string `json:"dateCreated"`
	MimeType     string `json:"mimeType"`
	LastModified string `json:"lastModified"`
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

func grpcFunc(c *Folder, requestType string, preferredCase string) error {

	path, err := FindProjectRoot("server.env")
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	env, err := loadEnvFile(path)
	if err != nil {
		log.Fatalf("failed to read env file: %v", err)
	}

	port := env["PYTHON_SERVER"]

	if requestType != "METADATA" && requestType != "CLUSTERING" && requestType != "KEYWORDS" {
		return fmt.Errorf("invalid requestType: %s", requestType)
	}

	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create new grpc go client")
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	shh, found := os.LookupEnv("SFM_SERVER_SECRET")
	if !found {
		fmt.Println("secret not found")
		return errors.New("server secret not found error")
	}

	req := &pb.DirectoryRequest{
		Root:          convertFolderToProto(*c),
		RequestType:   requestType,
		PreferredCase: preferredCase,
		ServerSecret:  shh,
	}
	// printDirectoryWithMetadata(req.Root, 0)

	resp, err := client.SendDirectoryStructure(ctx, req)

	if err != nil {
		log.Fatalf("SendDirectoryStructure RPC failed: %v", err)
	}

	// fmt.Printf("Server returned root directory: name=%q, path=%q/n", resp.Root.GetName(), resp.Root.GetPath())

	switch requestType {
	case "KEYWORDS":
		mergeKeywordsInPlaceFromProto(resp.Root, c)
	default: // "METADATA", "CLUSTERING"
		mergeProtoToFolder(resp.Root, c)
	}

	return nil
}

func mergeProtoToFolder(dir *pb.Directory, existing *Folder) {
	if dir == nil || existing == nil {
		fmt.Println("mergeProtoToFolder called with nil")
		return
	}
	existing.Files = existing.Files[:0]
	existing.Subfolders = existing.Subfolders[:0]

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

	existing.Files = existing.Files[:0]
	existing.Subfolders = existing.Subfolders[:0]

	for _, file := range dir.Files {
		existing.Files = append(existing.Files, &File{
			Name:     file.Name,
			Path:     file.OriginalPath,
			Tags:     tagsToStrings(file.Tags),
			Metadata: metadataConverter(file.Metadata),
			Locked:   file.IsLocked,
			NewPath:  file.NewPath,
			Keywords: file.Keywords,
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

func mergeKeywordsInPlaceFromProto(dir *pb.Directory, existing *Folder) {
	if dir == nil || existing == nil {
		return
	}

	// Index existing files by Path once.
	files := make(map[string]*File, 1024)
	var index func(*Folder)
	index = func(f *Folder) {
		if f == nil {
			return
		}
		for _, fl := range f.Files {
			if fl != nil && fl.Path != "" {
				files[fl.Path] = fl
			}
		}
		for _, sf := range f.Subfolders {
			index(sf)
		}
	}
	index(existing)

	// Walk response and append keywords to matched files.
	var walk func(*pb.Directory)
	walk = func(d *pb.Directory) {
		if d == nil {
			return
		}
		for _, pf := range d.Files {
			if pf == nil {
				continue
			}
			key := strings.TrimSpace(pf.GetOriginalPath())
			if key == "" {
				key = strings.TrimSpace(pf.GetNewPath())
			}
			if exf, ok := files[key]; ok {
				exf.Keywords = AppendUniqueKeywords(exf.Keywords, pf.Keywords)
				// Only keywords are updated for this request type.
			}
		}
		for _, sd := range d.Directories {
			walk(sd)
		}
	}
	walk(dir)
}

func tagsToStrings(tags []*pb.Tag) []string {
	var tagStrings []string

	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.Name)
	}
	return tagStrings
}

// converts string version of tags to *pb.Tag array
func stringsToTags(stringTags []string) []*pb.Tag {
	var tags []*pb.Tag

	for _, tag := range stringTags {
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

// func printDirectoryWithMetadata(dir *pb.Directory, num int) {

// 	space := strings.Repeat("  ", num)
// 	fmt.Println(space + "Directory: " + dir.Name)
// 	fmt.Println(space + "Path: " + dir.Path)

// 	for _, file := range dir.Files {
// 		fmt.Println(space + "File name: " + file.Name)
// 		fmt.Println("keywords: ")

// 		fmt.Println("----")
// 	}

// 	for _, dir := range dir.Directories {
// 		newNum := num + 1
// 		printDirectoryWithMetadata(dir, newNum)
// 	}

// }
func printDirectoryWithMetadata(dir *pb.Directory, leftPad int) {
	if dir == nil {
		fmt.Println("(nil directory)")
		return
	}
	prefix := strings.Repeat(" ", leftPad)
	fmt.Printf("%s%s\n", prefix, dirLabel(dir))
	printDirChildren(dir, prefix)
}

func printDirChildren(d *pb.Directory, prefix string) {
	if d == nil {
		return
	}

	// Copy and sort for stable output: directories first, then files.
	subdirs := append([]*pb.Directory(nil), d.Directories...)
	files := append([]*pb.File(nil), d.Files...)

	sort.Slice(subdirs, func(i, j int) bool {
		ni := subdirs[i].Name
		if ni == "" {
			ni = subdirs[i].Path
		}
		nj := subdirs[j].Name
		if nj == "" {
			nj = subdirs[j].Path
		}
		return strings.ToLower(ni) < strings.ToLower(nj)
	})

	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	if len(subdirs) == 0 && len(files) == 0 {
		fmt.Printf("%s└── (empty)\n", prefix)
		return
	}

	// Print directories
	for i, sd := range subdirs {
		isLast := i == len(subdirs)-1 && len(files) == 0
		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}
		fmt.Printf("%s%s%s\n", prefix, branch, dirLabel(sd))
		printDirChildren(sd, nextPrefix)
	}

	// Print files
	for j, f := range files {
		isLast := j == len(files)-1
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		fmt.Printf("%s%s%s\n", prefix, branch, fileLabel(f))
		fmt.Printf("%s%s%s%s\n", prefix, branch, fileLabel(f), "'s keyw:")

		fmt.Printf("%s%s locked is: %s\n", prefix, branch, strconv.FormatBool(f.IsLocked))

	}
}

func dirLabel(d *pb.Directory) string {
	name := d.Name
	if strings.TrimSpace(name) == "" {
		name = "<unnamed>"
	}

	b := strings.Builder{}
	b.Grow(64)
	b.WriteString("[DIR] ")
	b.WriteString(name)

	if d.Path != "" && d.Path != name {
		b.WriteString(" - ")
		b.WriteString(d.Path)
	}

	b.WriteString(" [dirs=")
	b.WriteString(fmt.Sprintf("%d", len(d.Directories)))
	b.WriteString(" files=")
	b.WriteString(fmt.Sprintf("%d", len(d.Files)))
	b.WriteString("]")

	if d.IsLocked {
		b.WriteString(" [locked]")
	}

	return b.String()
}

func fileLabel(f *pb.File) string {
	name := f.Name
	if strings.TrimSpace(name) == "" {
		name = "<unnamed>"
	}
	// Keep this simple and safe (no assumptions about File fields beyond Name).
	return "[FILE] " + name
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

			populateKeywordsFromStoredJsonFile(c)

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

			//reads json stored keywords and adds to the current composite

			// starts extracting keywords to cover new files and changes in files
			GoExtractKeywords(c)

			go pythonExtractKeywords(c)

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
			NewPath:  file.NewPath,
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

// func createDirectoryJSONStructure(folder *Folder) []FileNode {
// 	var nodes []FileNode

// 	for _, file := range folder.Files {

// 		md := extractMetadata(file.Metadata)
// 		tags := file.Tags
// 		// if len(tags) == 0 {
// 		// 	tags = []string{"none"}
// 		// }

// 		nodes = append(nodes, FileNode{
// 			Name:     file.Name,
// 			Path:     file.Path,
// 			IsFolder: false,
// 			Tags:     tags,
// 			Metadata: md,
// 			Locked:   file.Locked,
// 			NewPath:  file.NewPath, // include NewPath for moving files
// 		})
// 	}

// 	for _, sub := range folder.Subfolders {
// 		// recurse first
// 		childNodes := createDirectoryJSONStructure(sub)

// 		nodes = append(nodes, FileNode{
// 			Name:     sub.Name,
// 			Path:     sub.Path,
// 			IsFolder: true,
// 			Tags:     sub.Tags,
// 			Metadata: &Metadata{},
// 			Children: childNodes,
// 			Locked:   sub.Locked,
// 			NewPath:  sub.NewPath,
// 		})
// 	}

// 	return nodes
// }
