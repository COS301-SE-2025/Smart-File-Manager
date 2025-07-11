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

	root := &Folder{
		Name:         managerName,
		Path:         cleanPath,
		CreationDate: time.Now(),
	}

	// Recursively scan filesystem
	if err := exploreDown(root, cleanPath); err != nil {
		return nil, fmt.Errorf("error exploring folder %q: %w", cleanPath, err)
	}

	autoLockHiddenFolders(root)

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

func autoLockHiddenFolders(folder *Folder) {
	for _, sub := range folder.Subfolders {
		autoLockHiddenFolders(sub)

		// Check if subfolder is hidden
		if strings.HasPrefix(sub.Name, ".") {
			folder.Locked = true
			fmt.Printf("Auto-locked folder '%s' because it contains hidden folder '%s'\n", folder.Path, sub.Name)
			break // No need to check further
		}
	}
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
