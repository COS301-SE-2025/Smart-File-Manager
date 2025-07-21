package filesystem

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var root string

func moveDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	composites := GetComposites()
	if len(composites) == 0 {
		http.Error(w, "No managers found", http.StatusNotFound)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	compositeName := r.URL.Query().Get("name")

	for _, item := range composites {
		if item.Name == compositeName {
			//build the directory structure for the smart manager
			CreateDirectoryStructure(item)
			// fmt.Println("Directory structure created for composite:", compositeName)
			//Move the content of the smart manager
			moveContent(item)
			// fmt.Println("Content moved for composite:", compositeName)
			w.Write([]byte("true"))
			return
		}
	}
	fmt.Println("Smart manager not found: ", compositeName)
	w.Write([]byte("false"))
	curDir, _ := os.Getwd()
	fmt.Println("Current working directory:", curDir)

}

func moveContent(item *Folder) {
	//Move the files&folders according to new path after sorting
}

func CreateDirectoryStructure(item *Folder) {
	root = getPath()
	root = filepath.Join(root, "archives", item.Name)
	//call the recursive function to create the directory structure
	CreateDirectoryStructureRecursive(item)
	//change root back to original path
	os.Chdir("filesystem")
}
func CreateDirectoryStructureRecursive(item *Folder) {
	if item == nil {
		return
	}

	if len(item.Subfolders) == 0 {
		targetPath := filepath.Join(root, item.NewPath)
		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			panic(err)
		}
		return
	}
	for _, subfolder := range item.Subfolders {
		CreateDirectoryStructureRecursive(subfolder)
	}
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Search upward for "Smart-File-Manager"
	for {
		if filepath.Base(dir) == "Smart-File-Manager" {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // hit root
		}
		dir = parent
	}
	panic("could not find Smart-File-Manager root")
}
