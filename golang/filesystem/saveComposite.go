package filesystem

import (
	"encoding/json"
	"os"
	"path/filepath"
)

//uses load tree struct directoryTreeJson

func saveCompositeDetails(c *Folder) {
	children := compositeToJsonStorageFormat(c)

	root := DirectoryTreeJson{
		Name:     c.Name,
		IsFolder: true,
		RootPath: c.Path,
		Children: children,
	}

	saveCompositeDetailsToFile(root)
}

func compositeToJsonStorageFormat(folder *Folder) []FileNode {
	var nodes []FileNode

	for _, file := range folder.Files {
		tags := file.Tags

		nodes = append(nodes, FileNode{
			Name:     file.Name,
			Path:     file.Path,
			IsFolder: false,
			Tags:     tags,
			Locked:   file.Locked,
		})
	}

	for _, sub := range folder.Subfolders {
		// recurse first
		childNodes := compositeToJsonStorageFormat(sub)

		nodes = append(nodes, FileNode{
			Name:     sub.Name,
			Path:     sub.Path,
			IsFolder: true,
			Tags:     sub.Tags,
			Metadata: &Metadata{},
			Children: childNodes,
			Locked:   sub.Locked,
		})
	}

	return nodes
}

var compositeStoragePath = filepath.Join("storage", "composite.json")

// function that stores the composite
func saveCompositeDetailsToFile(comp DirectoryTreeJson) error {
	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(compositeStoragePath), 0755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(compositeStoragePath, out, 0644)
}
