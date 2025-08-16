package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// uses load tree struct directoryTreeJson

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
			Keywords: file.Keywords,
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
			Children: childNodes,
			Locked:   sub.Locked,
		})
	}

	return nodes
}

// function that stores the composite
func saveCompositeDetailsToFile(comp DirectoryTreeJson) error {
	var filePath = filepath.Join("storage", (comp.Name + ".json"))
	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	out, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, out, 0644)
}

func populateKeywordsFromStoredJsonFile(comp *Folder) {
	fmt.Println("populateKeywordsFromStoredJsonFile called")
	var filePath = filepath.Join("storage", (comp.Name + ".json"))
	// If the file doesn't exist yet, start with empty
	data, err := os.ReadFile(filePath)

	if os.IsNotExist(err) {
		return
	} else if err != nil {
		return
	}

	var structure DirectoryTreeJson

	//populates recs
	if err := json.Unmarshal(data, &structure); err != nil {
		fmt.Println("error in unmarshaling of json")
		return
	}
	mergeDirectoryTreeToComposite(comp, &structure)

}

func mergeDirectoryTreeToComposite(comp *Folder, directory *DirectoryTreeJson) {

	fmt.Println("mergeDirectoryTreeToComposite called")
	for _, node := range directory.Children {
		if !node.IsFolder {
			path := node.Path

			var compositeFile *File = comp.GetFile(path)
			if compositeFile != nil {
				compositeFile.Keywords = node.Keywords
			}

		} else {
			helperMergeDirectoryTreeToComposite(comp, &node)
		}
	}
}

func helperMergeDirectoryTreeToComposite(comp *Folder, fileNode *FileNode) {
	fmt.Println("helperMergeDirectoryTreeToComposite called")
	for _, node := range fileNode.Children {
		if !node.IsFolder {
			path := node.Path

			var compositeFile *File = comp.GetFile(path)
			if compositeFile != nil {
				compositeFile.Keywords = node.Keywords
			}

		} else {
			helperMergeDirectoryTreeToComposite(comp, &node)
		}
	}

}
