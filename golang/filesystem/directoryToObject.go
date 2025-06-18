package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConvertToObject builds a Folder tree from the given path (absolute or relative).
func ConvertToObject(managerName, folderPath string) (*Folder, error) {
	// Resolve to absolute path
	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %q: %w", folderPath, err)
	}

	root := &Folder{
		Name:         managerName,
		Path:         absPath,
		CreationDate: time.Now(),
	}

	// Recursively scan filesystem
	if err := exploreDown(root, absPath); err != nil {
		return nil, fmt.Errorf("error exploring folder %q: %w", absPath, err)
	}

	return root, nil
}

// exploreDown reads the directory at path and adds subfolders/files to folder
func exploreDown(folder *Folder, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			sub := &Folder{
				Name:         entry.Name(),
				Path:         fullPath,
				CreationDate: info.ModTime(),
			}
			folder.AddSubfolder(sub)
			// recurse into subfolder
			if err := exploreDown(sub, fullPath); err != nil {
				fmt.Printf("warning: cannot explore %s: %v\n", fullPath, err)
			}
		} else {
			file := &File{
				Name:     entry.Name(),
				Path:     fullPath,
				Metadata: []*MetadataEntry{},
				Tags:     []string{},
			}
			folder.AddFile(file)
		}
	}
	return nil
}
