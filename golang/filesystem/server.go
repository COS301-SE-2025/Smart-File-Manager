package filesystem

import (
	// "encoding/json"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	//array of smartfile managers
	Composites []*Folder
	mu         sync.Mutex
)

// savePortToEnv updates or creates the GO_PORT entry in server.env
func savePortToEnv(port int) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Navigate to project root
	projectRoot := filepath.Dir(cwd)
	envFilePath := filepath.Join(projectRoot, "server.env")

	// Try current directory first if the above doesn't work
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		envFilePath = filepath.Join(cwd, "server.env")
	}

	// If still not found, try going up two levels
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		projectRoot = filepath.Dir(filepath.Dir(cwd))
		envFilePath = filepath.Join(projectRoot, "server.env")
	}

	var lines []string
	goPortFound := false

	// Read existing file if it exists
	if file, err := os.Open(envFilePath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "GO_PORT=") {
				lines = append(lines, fmt.Sprintf("GO_PORT=%d", port))
				goPortFound = true
			} else {
				lines = append(lines, line)
			}
		}
	}

	// Add GO_PORT if not found
	if !goPortFound {
		lines = append(lines, fmt.Sprintf("GO_PORT=%d", port))
	}

	// Write back to file
	content := strings.Join(lines, "\n")
	err = os.WriteFile(envFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to server.env: %v", err)
	}

	return nil
}

// findAvailablePort finds an available port starting from a base port
func findAvailablePort(basePort int) (int, error) {
	for port := basePort; port < basePort+100; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", basePort, basePort+99)
}

func isPathContained(parentPath, childPath string) bool {
	parentPath = ConvertToWSLPath(filepath.Clean(parentPath))
    childPath = ConvertToWSLPath(filepath.Clean(childPath))


	// Convert to absolute paths for accurate comparison
	parentAbs, err := filepath.Abs(parentPath)
	if err != nil {
		return false
	}

	childAbs, err := filepath.Abs(childPath)
	if err != nil {
		return false
	}

	// Ensure paths end with separator for accurate prefix checking
	if !strings.HasSuffix(parentAbs, string(filepath.Separator)) {
		parentAbs += string(filepath.Separator)
	}
	if !strings.HasSuffix(childAbs, string(filepath.Separator)) {
		childAbs += string(filepath.Separator)
	}

	// Check if child path starts with parent path
	return strings.HasPrefix(childAbs, parentAbs)
}

func checkDirectoryConflicts(newPath string) (bool, string, error) {
	// fmt.Println(" Manager directory conflicts started");
	newPath = ConvertToWSLPath(filepath.Clean(newPath))

	for _, comp := range Composites {
		existingPathAbs := comp.Path

		fmt.Println("New Path: "+ newPath);
		fmt.Println("Old Path: "+ existingPathAbs);
		// Exact match
		if existingPathAbs == newPath {
			return true, fmt.Sprintf("Directory is already managed by '%s'", comp.Name), nil
		}

		// New path is inside existing manager
		if strings.HasPrefix(newPath+string(os.PathSeparator), existingPathAbs+string(os.PathSeparator)) {
			return true, fmt.Sprintf("Directory is already contained within existing manager '%s' at path '%s'", comp.Name, comp.Path), nil
		}

		// Existing manager is inside new path
		if strings.HasPrefix(existingPathAbs+string(os.PathSeparator), newPath+string(os.PathSeparator)) {
			return true, fmt.Sprintf("New directory would contain existing manager '%s' at path '%s'", comp.Name, comp.Path), nil
		}
	}

	return false, "", nil
}

func addCompositeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" Manager creation started");
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	// Check if manager name already exists
	for _, comp := range Composites {
		if comp.Name == managerName {
			http.Error(w, "A smart file manager with that name already exists", http.StatusBadRequest)
			return
		}
	}

	mu.Lock()
	defer mu.Unlock()

	// Check for directory conflicts
		fmt.Println(" Manager directory conflicts called");

	hasConflict, _, err := checkDirectoryConflicts(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking directory conflicts: %v", err), http.StatusInternalServerError)
		return
	}
			fmt.Println(" Manager directory conflicts ended");

	if hasConflict {
		w.Write([]byte("false"))
		return
	}

	// Proceed with adding the manager
	err = AddManager(managerName, filePath)
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	w.Write([]byte("true"))
}

func addTagHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	tag := r.URL.Query().Get("tag")

	convertedPath := ConvertToWSLPath(filePath)
	if filePath == "" || tag == "" {
		w.Write([]byte("false"))
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		item := c.GetFile(convertedPath)
		if item != nil {
			c.AddTagToFile(convertedPath, tag)
			c.Display(0)
			saveCompositeDetails(c)
			w.Write([]byte("true"))
			return
		}
	}

	// fmt.Println("Item not found for path:", convertedPath)
	w.Write([]byte("false"))
}

func removeTagHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	tag := r.URL.Query().Get("tag")

	convertedPath := ConvertToWSLPath(filePath)

	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		// Check file
		if file := c.GetFile(convertedPath); file != nil {
			if file.RemoveTag(tag) {
				// fmt.Printf("Removed tag '%s' from file: %s\n", tag, convertedPath)
				saveCompositeDetails(c)
				w.Write([]byte("true"))
				return
			}
		}
		// Csheck folder
		if folder := c.GetSubfolder(convertedPath); folder != nil {
			if folder.RemoveTag(tag) {
				// fmt.Printf("Removed tag '%s' from folder: %s\n", tag, convertedPath)
				w.Write([]byte("true"))
				return
			}
		}
	}

	// fmt.Println("Tag or item not found for path:", convertedPath)
	w.Write([]byte("false"))
}

// Locks a file or folder and its children
func lockHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	name := r.URL.Query().Get("name")
	if path == "" || name == "" {
		w.Write([]byte("Parameter missing"))
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		if c.Name == name {
			c.LockByPath(path)
			saveCompositeDetails(c)
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))

}

// Unlocks a file or folder and its children
func unlockHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	if path == "" || name == "" {
		w.Write([]byte("Parameter missing"))
		return
	}
	for _, c := range Composites {
		if c.Name == name {
			c.UnlockByPath(path)
			saveCompositeDetails(c)
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))

}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()
	if path == "" || name == "" {
		w.Write([]byte("Parameter missing"))
		return
	}
	for _, c := range Composites {
		if c.Name == name {
			err := os.RemoveAll(path)
			if err == nil {
				c.RemoveFile(path)
			} else {
				fmt.Println("Error removing folder:", path, "Error:", err)
			}
			children := GoSidecreateDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				RootPath: c.Path,
				Children: children,
			}
			// w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}
	}

	w.Write([]byte("false"))

}

func deleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()
	if path == "" || name == "" {
		w.Write([]byte("Parameter missing"))
		return
	}
	for _, c := range Composites {
		if c.Name == name {
			err := os.RemoveAll(path)
			if err == nil {
				c.RemoveSubfolder(path)
			} else {
				fmt.Println("Error removing folder:", path, "Error:", err)
			}
			children := GoSidecreateDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				RootPath: c.Path,
				Children: children,
			}
			// w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}
	}
	w.Write([]byte("false"))

}

func deleteManagerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()
	if name == "" {
		w.Write([]byte("Parameter missing"))
		return
	}
	for i, c := range Composites {
		if c.Name == name {
			// Delete folder
			// os.RemoveAll(c.Path)
			// Remove from list of managers
			Composites = append(Composites[:i], Composites[i+1:]...)
			// Remove from type storage
			delete(ObjectMap, c.Path)

			data, err := os.ReadFile(managersFilePath)
			var recs []ManagerRecord

			if err == nil {
				// File exists — update entry
				if err := json.Unmarshal(data, &recs); err != nil {
					fmt.Println("error in unmarshaling of json")
					panic(err)
				}
				// Remove the record with the matching name
				for j := range recs {
					if recs[j].Name == name {
						recs = append(recs[:j], recs[j+1:]...)
						break
					}
				}
			} else if os.IsNotExist(err) {
				// Nothing to do if file doesn't exist
			} else {
				panic(err)
			}

			out, err := json.MarshalIndent(recs, "", "  ")
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(managersFilePath, out, 0644); err != nil {
				panic(err)
			}
			//remove storage file
			deleteCompositeDetailsFile(c.Name)
			// fmt.Println("Deleted manager")
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("Manager not found"))
}

func secretMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("apiSecret")
		apiSecret, found := os.LookupEnv("SFM_API_SECRET")
		if !found {
			fmt.Println("api secret not found")
			http.Error(w, "Server secret not configured", http.StatusInternalServerError)
			return
		}

		if secret != apiSecret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HandleRequests() {
	// Find an available port starting from 51000
	port, err := findAvailablePort(51000)
	if err != nil {
		fmt.Printf("Error finding available port: %v\n", err)
		port = 51000 // Fallback to original port
	} else {
		// Save the found port to server.env
		if err := savePortToEnv(port); err != nil {
			fmt.Printf("Warning: Could not save port to server.env: %v\n", err)
		} else {
			fmt.Printf("Saved port %d to server.env\n", port)
		}
	}

	http.Handle("/addDirectory", secretMiddleware(http.HandlerFunc(addCompositeHandler)))

	http.Handle("/addTag", secretMiddleware(http.HandlerFunc(addTagHandler)))
	http.Handle("/removeTag", secretMiddleware(http.HandlerFunc(removeTagHandler)))

	http.Handle("/loadTreeData", secretMiddleware(http.HandlerFunc(loadTreeDataHandlerGoOnly)))

	http.Handle("/sortTree", secretMiddleware(http.HandlerFunc(sortTreeHandler)))
	http.Handle("/startUp", secretMiddleware(http.HandlerFunc(startUpHandler)))

	http.Handle("/lock", secretMiddleware(http.HandlerFunc(lockHandler)))
	http.Handle("/unlock", secretMiddleware(http.HandlerFunc(unlockHandler)))

	http.Handle("/search", secretMiddleware(http.HandlerFunc(SearchHandler)))

	http.Handle("/keywordSearch", secretMiddleware(http.HandlerFunc(KeywordSearchHadler)))
	http.Handle("/isKeywordSearchReady", secretMiddleware(http.HandlerFunc(IsKeywordSearchReadyHander)))

	http.Handle("/moveDirectory", secretMiddleware(http.HandlerFunc(moveDirectoryHandler)))

	http.Handle("/findDuplicateFiles", secretMiddleware(http.HandlerFunc(findDuplicateFilesHandler)))

	http.Handle("/bulkAddTag", secretMiddleware(http.HandlerFunc(BulkAddTagHandler)))
	http.Handle("/bulkRemoveTag", secretMiddleware(http.HandlerFunc(BulkRemoveTagHandler)))

	http.Handle("/deleteFile", secretMiddleware(http.HandlerFunc(deleteFileHandler)))
	http.Handle("/deleteFolder", secretMiddleware(http.HandlerFunc(deleteFolderHandler)))
	http.Handle("/bulkDeleteFolders", secretMiddleware(http.HandlerFunc(BulkDeleteFolderHandler)))
	http.Handle("/bulkDeleteFiles", secretMiddleware(http.HandlerFunc(BulkDeleteFileHandler)))
	http.Handle("/deleteManager", secretMiddleware(http.HandlerFunc(deleteManagerHandler)))

	http.Handle("/returnType", secretMiddleware(http.HandlerFunc(ReturnTypeHandler)))

	http.Handle("/returnStats", secretMiddleware(http.HandlerFunc(StatHandler)))

	http.Handle("/setPreferredCase", secretMiddleware(http.HandlerFunc(SetPreferredCase)))

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}

}

// Getter for main.go
func GetComposites() []*Folder {
	mu.Lock()
	defer mu.Unlock()
	return Composites
}
