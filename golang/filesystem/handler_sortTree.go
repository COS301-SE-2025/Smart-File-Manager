package filesystem

import (
	"encoding/json"
	"log"
	"net/http"
	//grpc imports
)

// actual api endpoint function
func sortTreeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		if c.Name == name {
			// build the nested []FileNode
			err := grpcFunc(c, "CLUSTERING")
			if err != nil {
				log.Fatalf("grpcFunc failed: %v", err)
				http.Error(w, "internal server error, GRPC CALLED WRONG", http.StatusInternalServerError)
			}
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
