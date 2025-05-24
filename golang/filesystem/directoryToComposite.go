package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func convertToComposite(managerID string, managerName string, filePath string) *Folder {

	fmt.Println("Converting: ", managerName, " to composite")
	//creating managedItem
	root := &Folder{managedItem: managedItem{
		itemID:       managerID,
		itemName:     managerName,
		itemPath:     filePath,
		creationDate: time.Now(),
	}}

	// Recursively populate the folder with its contents
	err := exploreDown(root, filePath)
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
