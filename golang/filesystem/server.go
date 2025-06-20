package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"
)

var (
	//array of smartfile managers
	composites []*Folder
	mu         sync.Mutex
)

func getCompositeHandler(w http.ResponseWriter, r *http.Request) {
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")
	fmt.Println("PATH", filePath)
	composite, err := ConvertToObject(managerName, filePath)
	if err != nil || composite == nil {
		w.Write([]byte("false"))
		return
	}

	mu.Lock()
	composites = append(composites, composite)
	mu.Unlock()

	fmt.Println("Composite added:")
	composite.Display(0)
	w.Write([]byte("true"))
}

func removeCompositeHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	convertedPath := ConvertToWSLPath(filePath)

	mu.Lock()
	for i, item := range composites {
		if item.Path == convertedPath {
			composites = slices.Delete(composites, i, i+1)
			break
		}
	}
	mu.Unlock()

	fmt.Println("Composite removed:", convertedPath)
	w.Write([]byte("true"))
}

func addTagHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	tag := r.URL.Query().Get("tag")

	convertedPath := ConvertToWSLPath(filePath)

	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		item := c.GetFile(convertedPath)
		if item != nil {
			c.AddTagToFile(convertedPath, tag)
			c.Display(0)
			w.Write([]byte("true"))
			return
		}
	}

	fmt.Println("Item not found for path:", convertedPath)
	w.Write([]byte("false"))
}

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
			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

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

func HandleRequests() {

	// path, _ := os.Getwd()
	// fmt.Println("THE PATH: " + path)
	// path = filepath.Dir(path)
	// path = filepath.Join(path, "python/testing")
	// fmt.Println("THE PATH: " + path)

	http.HandleFunc("/addDirectory", getCompositeHandler)
	http.HandleFunc("/removeDirectory", removeCompositeHandler)
	http.HandleFunc("/addTag", addTagHandler)
	http.HandleFunc("/loadTreeData", loadTreeDataHandler)
	fmt.Println("Server started on port 51000")
	// http.ListenAndServe(":51000", nil)
	http.ListenAndServe("0.0.0.0:51000", nil)

}

// Getter for main.go
func GetComposites() []*Folder {
	mu.Lock()
	defer mu.Unlock()
	return composites
}
