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

// uses temp files to prevent races / overwritting a file that is being read
func saveCompositeDetailsToFile(comp DirectoryTreeJson) error {
	filePath := filepath.Join("storage", comp.Name+".json")
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	out, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "tmp-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	// Clean up temp file on any error.
	cleanup := func(e error) error {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return e
	}

	if _, err := tmp.Write(out); err != nil {
		return cleanup(err)
	}

	if err := tmp.Sync(); err != nil {
		return cleanup(err)
	}

	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	if err := os.Rename(tmpName, filePath); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}

	return nil
}

func populateKeywordsFromStoredJsonFile(comp *Folder) {
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
				compositeFile.Tags = node.Tags
				compositeFile.Locked = node.Locked
			}

		} else {
			helperMergeDirectoryTreeToComposite(comp, &node)
		}
	}
}

func helperMergeDirectoryTreeToComposite(comp *Folder, fileNode *FileNode) {
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

func deleteCompositeDetailsFile(compName string) error {
	filePath := filepath.Join("storage", compName+".json")
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // already gone
		}
		return err
	}
	return nil
}
