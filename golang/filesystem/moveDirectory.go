package filesystem

import (
	"fmt"
	"net/http"
)

func moveDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("moveDirectory called")
	// w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		if c.Name == name {
			// Create new file structure
			createDirectoryStructure()
			// Move the content to the new structure
			moveContent()
			fmt.Println("Directory moved successfully for:", name)
			w.Write([]byte("true"))
			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

func moveContent() {
	//Move the /files/folders according to new path after sorting
}

func createDirectoryStructure() {
	//Create the directory structure for the sorted content

}
