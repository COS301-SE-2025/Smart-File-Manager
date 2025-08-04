package filesystem

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
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
	// fmt.Println("PATH", filePath)
	composite, err := ConvertToObject(managerName, filePath)
	if err != nil || composite == nil {
		w.Write([]byte("false"))
		return
	}

	mu.Lock()
	// Composites = append(Composites, composite)
	//appendng happens in this:
	AddManager(managerName, filePath)
	mu.Unlock()

	// fmt.Println("Composite added:")
	// composite.Display(0)
	w.Write([]byte("true"))
}

func removeCompositeHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	convertedPath := ConvertToWSLPath(filePath)

	mu.Lock()
	for i, item := range Composites {
		if item.Path == convertedPath {
			Composites = slices.Delete(Composites, i, i+1)
			break
		}
	}
	mu.Unlock()

	// fmt.Println("Composite removed:", convertedPath)
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
			fmt.Println("LOCKED FILE")
			w.Write([]byte("true"))
			return
		} else {
			w.Write([]byte("false"))
			return
		}
	}

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
			w.Write([]byte("true"))
			return
		} else {
			w.Write([]byte("false"))
			return
		}
	}

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
			err := os.Remove(path)
			if err != nil {
				panic(err)
			}
			c.RemoveFile(path)
			// c.Display(0)
			w.Write([]byte("true"))
			return
		} else {
			w.Write([]byte("false"))
			return
		}
	}
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
			if err != nil {
				panic(err)
			}
			c.RemoveSubfolder(path)
			// c.Display(0)
			w.Write([]byte("true"))
			return
		} else {
			w.Write([]byte("false"))
			return
		}
	}
}

func HandleRequests() {

	// path, _ := os.Getwd()
	// fmt.Println("THE PATH: " + path)
	// path = filepath.Dir(path)
	// path = filepath.Join(path, "python/testing")
	// fmt.Println("THE PATH: " + path)

	http.HandleFunc("/addDirectory", addCompositeHandler)
	http.HandleFunc("/removeDirectory", removeCompositeHandler)

	http.HandleFunc("/addTag", addTagHandler)
	http.HandleFunc("/removeTag", removeTagHandler)

	http.HandleFunc("/loadTreeData", loadTreeDataHandlerGoOnly)
	// http.HandleFunc("/loadTreeData", loadTreeDataHandler)
	http.HandleFunc("/sortTree", sortTreeHandler)
	http.HandleFunc("/startUp", startUpHandler)

	http.HandleFunc("/lock", lockHandler)
	http.HandleFunc("/unlock", unlockHandler)

	http.HandleFunc("/search", SearchHandler)

	http.HandleFunc("/moveDirectory", moveDirectoryHandler)

	http.HandleFunc("/findDuplicateFiles", findDuplicateFilesHandler)

	http.HandleFunc("/bulkAddTag", BulkAddTagHandler)
	http.HandleFunc("/bulkRemoveTag", BulkRemoveTagHandler)

	http.HandleFunc("/deleteFile", deleteFileHandler)
	http.HandleFunc("/deleteFolder", deleteFolderHandler)
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
