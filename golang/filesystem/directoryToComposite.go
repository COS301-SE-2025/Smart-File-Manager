package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func convertToComposite(managerID string, managerName string, filePath string) *Folder {

	fmt.Println("Converting: ", managerName, " to composite")
	var newPath = filePath
	if IsWindowsPath(filePath) {
		newPath = ConvertWindowsToWSLPath(filePath)
	}

	//creating managedItem
	root := &Folder{managedItem: managedItem{
		itemID:       managerID,
		itemName:     managerName,
		itemPath:     newPath,
		creationDate: time.Now(),
	}}

	// Recursively populate the folder with its contents
	err := exploreDown(root, newPath)
	if err != nil {
		fmt.Println("Error exploring folder:", err)
	}

	return root
}

// Recursively builds a composite tree from the directory structure
func exploreDown(folder *Folder, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		info, err := entry.Info()
		if err != nil {
			fmt.Println("Failed to get info for:", fullPath)
			continue
		}

		if entry.IsDir() {
			fmt.Println("Found folder:", fullPath)
			subFolder := &Folder{
				managedItem: managedItem{
					itemName:     entry.Name(),
					itemPath:     fullPath,
					creationDate: info.ModTime(),
				},
			}
			folder.AddItem(subFolder)
			err := exploreDown(subFolder, fullPath)
			if err != nil {
				fmt.Println("Error exploring subfolder:", err)
			}
		} else {
			fmt.Println("Found file:", fullPath)
			file := &File{
				managedItem: managedItem{
					itemName:     entry.Name(),
					itemPath:     fullPath,
					creationDate: info.ModTime(),
					fileType:     detectFileType(info), // optional
				},
			}
			err := folder.AddItem(file)
			if err != nil {
				fmt.Println("Error adding file:", err)
			}
		}
	}

	return nil
}
func detectFileType(info fs.FileInfo) string {
	return filepath.Ext(info.Name())
}

// IsWindowsPath returns true if the path looks like a Windows absolute path
func IsWindowsPath(path string) bool {
	return len(path) > 2 &&
		path[1] == ':' &&
		(path[2] == '\\' || path[2] == '/')
}

func ConvertWindowsToWSLPath(winPath string) string {
	winPath = strings.ReplaceAll(winPath, "\\", "/")
	if len(winPath) > 2 && winPath[1] == ':' {
		drive := strings.ToLower(string(winPath[0]))
		return "/mnt/" + drive + winPath[2:]
	}
	return winPath
}
func deleteComposite(f **Folder) { //will need to be modified when we start storing composites
	if f == nil || *f == nil {
		fmt.Println("Nothing to delete.")
		return
	}
	*f = nil
	fmt.Println("Composite deleted successfully.")
}
