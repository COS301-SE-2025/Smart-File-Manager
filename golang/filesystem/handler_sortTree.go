package filesystem

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	//grpc imports
)

// actual api endpoint function
func sortTreeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	caseParam := r.URL.Query().Get("case")
	if caseParam == "" {
		caseParam = "CAMEL"
	}
	allowedCases := []string{"CAMEL", "SNAKE", "PASCAL", "KEBAB"}
	if !slices.Contains(allowedCases, caseParam) {
		http.Error(w, "Invalid case type.", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		if c.Name == name {
			// build the nested []FileNode
			err := grpcFunc(c, "CLUSTERING", caseParam)
			if err != nil {
				log.Fatalf("grpcFunc failed: %v", err)
				http.Error(w, "internal server error, GRPC CALLED WRONG", http.StatusInternalServerError)
			}

			// PrettyPrintFolder(c, "")

			children := GoSidecreateDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				RootPath: c.Path,
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
