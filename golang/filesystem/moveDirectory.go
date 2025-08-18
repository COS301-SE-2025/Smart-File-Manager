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
			//update composite in memory
			newObj, err := ConvertToObject(item.Name, item.Path)
			if err != nil {
				log.Printf("Error converting to object: %v", err)
			}
			*item = *newObj
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

	item.Path = filepath.Join(root, item.Name)
	if originalPath != item.Path {
		os.RemoveAll(originalPath)

	}

	// Path to managers storage file
	managersFilePath := filepath.Join(getPath(), "golang", managersFilePath)
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
		targetPath := filepath.Join(root, file.NewPath)

		// Create target directory if it doesn't exist
		targetDir := filepath.Dir(targetPath)
		os.MkdirAll(targetDir, os.ModePerm)

		// Handle duplicate files by generating unique names
		finalTargetPath := generateUniqueFilePath(targetPath)

		if err := os.Rename(sourcePath, finalTargetPath); err != nil {
			log.Printf("Error moving file %s to %s: %v", sourcePath, finalTargetPath, err)
		}
	}

	for _, subfolder := range item.Subfolders {
		moveContentRecursive(subfolder)
	}
}

func generateUniqueFilePath(targetPath string) string {
	// If file doesn't exist, return original path
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return targetPath
	}

	// File exists, generate unique name
	dir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	counter := 1
	for {
		newFilename := fmt.Sprintf("%s_(%d)%s", nameWithoutExt, counter, ext)
		newPath := filepath.Join(dir, newFilename)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
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
