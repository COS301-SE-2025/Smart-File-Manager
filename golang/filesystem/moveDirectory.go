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

	for i, item := range Composites {
		// fmt.Printf("Checking manager: %s\n", item.Name)
		if item.Name == compositeName {
			// fmt.Printf("found manager: %s\n", item.Name)
			CreateDirectoryStructure(item)
			moveContent(item)

			name := item.Name
			path := item.Path

			Composites = append(Composites[:i], Composites[i+1:]...)
			delete(ObjectMap, item.Path)

			data, err := os.ReadFile(managersFilePath)
			var recs []ManagerRecord

			if err == nil {
				if err := json.Unmarshal(data, &recs); err != nil {
					fmt.Println("error in unmarshaling of json")
					panic(err)
				}
				for j := range recs {
					if recs[j].Name == name {
						recs = append(recs[:j], recs[j+1:]...)
						break
					}
				}
			} else if os.IsNotExist(err) {

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

			err = AddManager(name, path)
			if err != nil {
				log.Printf("Error adding manager: %v", err)
			}

			err = UpdateStoredPathsFromComposite(item)
			if err != nil {
				log.Printf("Error updating stored paths from composite: %v", err)
			}
			fmt.Println("responding with true")

			w.Write([]byte("true"))
			return
		}
	}
	fmt.Println("Smart manager not found: ", compositeName)
	w.Write([]byte("false"))
}

func moveContent(item *Folder) {
	parentDir := filepath.Dir(item.Path)
	originalPath := item.Path

	root = parentDir

	if err := os.MkdirAll(root, 0755); err != nil {
		panic(err)
	}

	moveContentRecursive(item)

	item.Path = filepath.Join(root, item.Name)
	if originalPath != item.Path {
		os.RemoveAll(originalPath)
	}

	managersFilePath := filepath.Join(getPath(), managersFilePath)

	data, err := os.ReadFile(managersFilePath)
	var recs []ManagerRecord

	if err == nil {
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

		targetDir := filepath.Dir(targetPath)
		os.MkdirAll(targetDir, os.ModePerm)

		finalTargetPath := generateUniqueFilePath(targetPath)

		if err := os.Rename(sourcePath, finalTargetPath); err != nil {
			log.Printf("Error moving file %s to %s: %v", sourcePath, finalTargetPath, err)
		} else {
			file.Path = finalTargetPath
		}
	}

	for _, subfolder := range item.Subfolders {
		subfolder.Path = filepath.Join(root, subfolder.NewPath)
		moveContentRecursive(subfolder)
	}
}

func generateUniqueFilePath(targetPath string) string {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return targetPath
	}

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
		subfolder.Path = filepath.Join(root, subfolder.NewPath)
		CreateDirectoryStructureRecursive(subfolder)
	}
}

func UpdateStoredPathsFromComposite(comp *Folder) error {
	filePath := filepath.Join("storage", comp.Name+".json")
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		saveCompositeDetails(comp)
		return nil
	} else if err != nil {
		return err
	}

	var oldStructure DirectoryTreeJson
	if err := json.Unmarshal(data, &oldStructure); err != nil {
		return err
	}

	newStructure := DirectoryTreeJson{
		Name:     comp.Name,
		IsFolder: true,
		RootPath: comp.Path,
		Children: buildNodesWithPreservedMetadata(comp, &oldStructure),
	}

	return saveCompositeDetailsToFile(newStructure)
}

func buildNodesWithPreservedMetadata(folder *Folder, oldStructure *DirectoryTreeJson) []FileNode {
	var nodes []FileNode

	oldPathMap := make(map[string]FileNode)
	buildPathMap(oldStructure.Children, oldPathMap)

	for _, file := range folder.Files {
		node := FileNode{
			Name:     file.Name,
			Path:     file.Path,
			IsFolder: false,
			Keywords: file.Keywords,
			Tags:     file.Tags,
			Locked:   file.Locked,
		}

		if oldNode, exists := findNodeByName(oldPathMap, file.Name, false); exists {
			if len(node.Keywords) == 0 {
				node.Keywords = oldNode.Keywords
			}
			if len(node.Tags) == 0 {
				node.Tags = oldNode.Tags
			}
			if !node.Locked {
				node.Locked = oldNode.Locked
			}
		}

		nodes = append(nodes, node)
	}

	for _, sub := range folder.Subfolders {
		childNodes := buildNodesWithPreservedMetadata(sub, oldStructure)

		node := FileNode{
			Name:     sub.Name,
			Path:     sub.Path,
			IsFolder: true,
			Tags:     sub.Tags,
			Children: childNodes,
			Locked:   sub.Locked,
		}

		if oldNode, exists := findNodeByName(oldPathMap, sub.Name, true); exists {
			if len(node.Tags) == 0 {
				node.Tags = oldNode.Tags
			}
			if !node.Locked {
				node.Locked = oldNode.Locked
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func buildPathMap(nodes []FileNode, pathMap map[string]FileNode) {
	for _, node := range nodes {
		pathMap[node.Path] = node
		pathMap[node.Name] = node

		if node.IsFolder {
			buildPathMap(node.Children, pathMap)
		}
	}
}

func findNodeByName(pathMap map[string]FileNode, name string, isFolder bool) (FileNode, bool) {
	if node, exists := pathMap[name]; exists && node.IsFolder == isFolder {
		return node, true
	}
	return FileNode{}, false
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func cleanManagerPrefix(path, managerName string) string {
	parts := strings.Split(path, string(os.PathSeparator))

	cleaned := []string{}
	seenManager := false
	for _, p := range parts {
		if p == managerName {
			if seenManager {
				continue
			}
			seenManager = true
		}
		cleaned = append(cleaned, p)
	}

	result := filepath.Join(cleaned...)
	// Preserve leading slash for absolute paths
	if filepath.IsAbs(path) && !strings.HasPrefix(result, string(os.PathSeparator)) {
		result = string(os.PathSeparator) + result
	}
	return result
}
