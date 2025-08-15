package filesystem

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
			// Build the directory structure in original path
			CreateDirectoryStructure(item)
			// Move the content into new manager folder
			moveContent(item)
			w.Write([]byte("true"))
			return
		}
	}
	fmt.Println("Smart manager not found: ", compositeName)
	w.Write([]byte("false"))
}

func moveContent(item *Folder) {
	// Get parent directory of original path (e.g. "/home/...")
	parentDir := filepath.Dir(item.Path)
	originalPath := item.Path

	// New root = parent directory + manager name
	root = parentDir

	// Create new root folder if not exists
	if err := os.MkdirAll(root, 0755); err != nil {
		panic(err)
	}

	moveContentRecursive(item)
	os.RemoveAll(originalPath)
	item.Path = filepath.Join(root, item.Name)

	// Path to managers storage file
	managersFilePath := filepath.Join(getPath(), "golang", "storage", "main.json")
	data, err := os.ReadFile(managersFilePath)
	var recs []ManagerRecord

	if err == nil {
		// File exists — update entry
		if err := json.Unmarshal(data, &recs); err != nil {
			fmt.Println("error in unmarshaling of json")
			panic(err)
		}
		for i := range recs {
			if recs[i].Name == item.Name {
				recs[i].Path = item.Path
			}
		}
	} else if os.IsNotExist(err) {
		// Create storage folder if missing
		if err := os.MkdirAll(filepath.Dir(managersFilePath), 0755); err != nil {
			panic(err)
		}
		recs = append(recs, ManagerRecord{item.Name, item.Path})
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

}

func moveContentRecursive(item *Folder) {
	if item == nil {
		return
	}

	for _, file := range item.Files {
		sourcePath := file.Path

		// Clean the new path to avoid triple manager names
		cleanedNewPath := cleanManagerPrefix(file.NewPath, item.Name)
		targetPath := filepath.Join(root, cleanedNewPath)

		// Also update the composite
		file.NewPath = cleanedNewPath

		os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err := os.Rename(sourcePath, targetPath); err != nil {
			log.Printf("Error moving file %s: %v", sourcePath, err)
		}
	}
	for _, subfolder := range item.Subfolders {
		moveContentRecursive(subfolder)
	}
}

func CreateDirectoryStructure(item *Folder) {
	// New root: original directory + manager name
	root = filepath.Join(item.Path, item.Name)
	if err := os.MkdirAll(root, 0755); err != nil {
		panic(err)
	}
	CreateDirectoryStructureRecursive(item)
}

func CreateDirectoryStructureRecursive(item *Folder) {
	if item == nil {
		return
	}

	if len(item.Subfolders) == 0 {
		targetPath := filepath.Join(root, item.NewPath)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
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

func cleanManagerPrefix(path, managerName string) string {
	parts := strings.Split(path, string(os.PathSeparator))

	// Remove repeated managerName prefixes
	cleaned := []string{}
	seenManager := false
	for _, p := range parts {
		if p == managerName {
			if seenManager {
				continue // skip extra occurrences
			}
			seenManager = true
		}
		cleaned = append(cleaned, p)
	}

	return filepath.Join(cleaned...)
}
