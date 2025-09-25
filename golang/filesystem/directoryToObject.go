package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConvertToObject builds a Folder tree from the given path (relative, WSL, or converted from Windows).
func ConvertToObject(managerName, folderPath string) (*Folder, error) {
	// Convert Windows path to WSL format if needed
	cleanPath := ConvertToWSLPath(folderPath)
	// cleanPath := folderPath

	root := &Folder{
		Name:         managerName,
		Path:         cleanPath,
		HasKeywords:  false,
		CreationDate: time.Now(),
	}

	// Recursively scan filesystem
	if err := exploreDown(root, cleanPath); err != nil {
		return nil, fmt.Errorf("error exploring folder %q: %w", cleanPath, err)
	}

	return root, nil
}

// exploreDown reads the directory at path and adds subfolders/files to folder
// It automatically locks the folder and all its descendants if it contains a hidden subfolder.
func exploreDown(folder *Folder, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			sub := &Folder{
				Name:         name,
				Path:         fullPath,
				CreationDate: info.ModTime(),
				Locked:       false,
			}
			folder.AddSubfolder(sub)
			if err := exploreDown(sub, fullPath); err != nil {
				// fmt.Printf("warning: cannot explore %s: %v\n", fullPath, err)
			}
		} else {
			file := &File{
				Name:     name,
				Path:     fullPath,
				Metadata: []*MetadataEntry{},
				Tags:     []string{},
				Locked:   false,
			}
			folder.AddFile(file)
		}
	}

	for _, sub := range folder.Subfolders {
		if strings.HasPrefix(sub.Name, ".") {
			folder.LockByPath(folder.Path)
			folder.Locked = false
			// fmt.Printf("Auto-locked folder '%s' and contents because it contains hidden folder '%s'\n", folder.Path, sub.Name)
			break
		}
	}
	for _, file := range folder.Files {
		if strings.HasPrefix(file.Name, ".") {
			file.Lock()
			// fmt.Printf("Auto-locked folder '%s' and contents because it contains hidden folder '%s'\n", folder.Path, sub.Name)
			break
		} else if strings.HasPrefix(file.Name, "~") {
			file.Lock()
			// fmt.Printf("Auto-locked folder '%s' and contents because it contains hidden folder '%s'\n", folder.Path, sub.Name)
			break
		}
	}
	return nil
}

func ConvertToWSLPath(winPath string) string {
	winPath = strings.TrimSpace(winPath)
	winPath = strings.ReplaceAll(winPath, "\\", "/")

	if len(winPath) > 2 && winPath[1] == ':' {
		drive := strings.ToLower(string(winPath[0]))
		rest := winPath[2:]
		return "/mnt/" + drive + rest
	}

	return winPath
}
