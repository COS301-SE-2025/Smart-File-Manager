package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var root string

func moveDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	compositeName := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, item := range Composites {
		fmt.Printf("Checking manager: %s\n", item.Name)
		if item.Name == compositeName {
			fmt.Printf("found manager: %s\n", item.Name)
			//build the directory structure for the smart manager
			CreateDirectoryStructure(item)
			//move the content of the smart manager to the new path
			moveContent(item)
			w.Write([]byte("true"))
			return
		}
	}
	fmt.Println("Smart managerpoop not found: ", compositeName)
	w.Write([]byte("false"))

}

func moveContent(item *Folder) {
	//Move the files&folders according to new path after sorting
	root = getPath()
	root = filepath.Join(root, "archives", item.Name)
	moveContentRecursive(item)
	//change root back to original path
	item.Path = root
	managersFilePath := getPath()
	managersFilePath = filepath.Join(managersFilePath, "golang", "storage", "main.json")
	//read storage
	data, err := os.ReadFile(managersFilePath)
	var exist bool
	if os.IsNotExist(err) {
		exist = false
	}

	var recs []ManagerRecord

	//populates recs
	if exist {
		if err := json.Unmarshal(data, &recs); err != nil {
			fmt.Println("error in unmarshaling of json")
			panic(err)
		}

		for i := range recs {
			if recs[i].Name == item.Name {
				recs[i].Path = item.Path
			}
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(managersFilePath), 0755); err != nil {
			panic(err)
		}
		record := ManagerRecord{
			item.Name,
			item.Path,
		}
		recs = append(recs, record)
	}

	out, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(managersFilePath, out, 0644)
	if err != nil {
		panic(err)
	}

	os.Chdir("filesystem")
}

func moveContentRecursive(item *Folder) {
	if item == nil {
		return
	}

	for _, file := range item.Files {
		sourcePath := file.Path
		targetPath := filepath.Join(root, file.NewPath)
		// Move the file
		err := os.Rename(sourcePath, targetPath)
		if err != nil {
			panic(err)
		}
		file.Path = targetPath // Update the file's path to the new location

	}
	for _, subfolder := range item.Subfolders {
		moveContentRecursive(subfolder)
	}
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
		targetPath := filepath.Join(root, item.Path)
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

// helper functions
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
