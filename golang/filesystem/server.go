package filesystem

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	//array of smartfile managers
	Composites []*Folder
	mu         sync.Mutex
)

func addCompositeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("addDirectory called")
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	for _, comp := range Composites {
		if comp.Name == managerName {
			http.Error(w, "A smart file manager with that name already exists", http.StatusBadRequest)
			return
		}
	}

	mu.Lock()
	// Composites = append(Composites, composite)
	//appendng happens in this:
	err := AddManager(managerName, filePath)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	mu.Unlock()

	w.Write([]byte("true"))
}

func addTagHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	tag := r.URL.Query().Get("tag")

	convertedPath := ConvertToWSLPath(filePath)

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
			w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusOK)
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
			fmt.Println("Deleted manager")
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

	// path, _ := os.Getwd()
	// fmt.Println("THE PATH: " + path)
	// path = filepath.Dir(path)
	// path = filepath.Join(path, "python/testing")
	// fmt.Println("THE PATH: " + path)

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

	fmt.Println("Server started on port 51000")

	// http.ListenAndServe(":51000", nil)
	http.ListenAndServe("0.0.0.0:51000", nil)

}

// Getter for main.go
func GetComposites() []*Folder {
	mu.Lock()
	defer mu.Unlock()
	return Composites
}
