package filesystem

import (
	"fmt"
	"net/http"
	"slices"
	"sync"
)

var (
	composites []*Folder
	mu         sync.Mutex
)

func getCompositeHandler(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("id")
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	composite := ConvertToComposite(managerID, managerName, filePath)
	if composite == nil {
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
	convertedPath := ConvertWindowsToWSLPath(filePath)

	mu.Lock()
	for i, item := range composites {
		if item.GetPath() == convertedPath {
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
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	convertedPath := ConvertWindowsToWSLPath(filePath)

	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		item := c.GetItem(convertedPath)
		if item != nil {
			item.AddTag(key, value)
			fmt.Println("Tag added to", convertedPath, ":", key, "=", value)
			c.Display(0)
			w.Write([]byte("true"))
			return
		}
	}

	fmt.Println("Item not found for path:", convertedPath)
	w.Write([]byte("false"))
}

func HandleRequests() {
	http.HandleFunc("/addDirectory", getCompositeHandler)
	http.HandleFunc("/removeDirectory", removeCompositeHandler)
	http.HandleFunc("/addTag", addTagHandler)
	http.ListenAndServe(":51000", nil)
}

// Getter for main.go
func GetComposites() []*Folder {
	mu.Lock()
	defer mu.Unlock()
	return composites
}
