package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
)

// uses load tree struct directoryTreeJson

func saveCompositeDetails(c *Folder) {

	if c == nil {
		return
	}
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
	if folder == nil {
		return nil
	}
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
	for _, node := range directory.Children {
		if !node.IsFolder {
			fmt.Println(node.Name)
			path := node.Path

			compositeFile := comp.GetFile(path)
			if compositeFile != nil {
				compositeFile.Keywords = node.Keywords
				compositeFile.Tags = node.Tags
				compositeFile.Locked = node.Locked
			}

			// fmt.Println("is locked for file : " + compositeFile.Name)
			// fmt.Println(strconv.FormatBool(compositeFile.Locked))
			// fmt.Println("\n ====")

		} else {
			helperMergeDirectoryTreeToComposite(comp, &node)
		}
	}
	fmt.Println("======END OFmergeDirectoryTreeToComposite========")

	// PrettyPrintFolder(comp, "")

}

func helperMergeDirectoryTreeToComposite(comp *Folder, fileNode *FileNode) {
	for _, node := range fileNode.Children {
		if !node.IsFolder {
			path := node.Path

			compositeFile := comp.GetFile(path)
			if compositeFile != nil {
				compositeFile.Keywords = node.Keywords
				compositeFile.Tags = node.Tags
				compositeFile.Locked = node.Locked
			}
			// fmt.Println("path : " + compositeFile.Path)

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

func printDirTreeKwTagsLock(root *DirectoryTreeJson, leftPad int) {
	if root == nil {
		fmt.Println("(nil tree)")
		return
	}
	prefix := strings.Repeat(" ", leftPad)
	fmt.Printf("%s[ROOT] %s\n", prefix, safeName(root.Name))
	printFileNodeChildren(root.Children, prefix)
}

func printFileNodeChildren(nodes []FileNode, prefix string) {
	if len(nodes) == 0 {
		fmt.Printf("%s└── (empty)\n", prefix)
		return
	}

	// Directories first, then files; stable, case-insensitive by name.
	dirs := make([]FileNode, 0, len(nodes))
	files := make([]FileNode, 0, len(nodes))
	for _, n := range nodes {
		if n.IsFolder {
			dirs = append(dirs, n)
		} else {
			files = append(files, n)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		ni := strings.ToLower(strings.TrimSpace(dirs[i].Name))
		nj := strings.ToLower(strings.TrimSpace(dirs[j].Name))
		return ni < nj
	})
	sort.Slice(files, func(i, j int) bool {
		ni := strings.ToLower(strings.TrimSpace(files[i].Name))
		nj := strings.ToLower(strings.TrimSpace(files[j].Name))
		return ni < nj
	})

	// Merge into single ordered list for correct tree branches.
	total := len(dirs) + len(files)
	entries := make([]FileNode, 0, total)
	entries = append(entries, dirs...)
	entries = append(entries, files...)

	for i, n := range entries {
		isLast := i == total-1
		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		fmt.Printf("%s%s%s\n", prefix, branch, nodeLabel(n))
		if n.IsFolder {
			printFileNodeChildren(n.Children, nextPrefix)
		}
	}
}

func nodeLabel(n FileNode) string {
	b := strings.Builder{}
	b.Grow(128)

	if n.IsFolder {
		b.WriteString("[DIR] ")
	} else {
		b.WriteString("[FILE] ")
	}
	b.WriteString(safeName(n.Name))

	if len(n.Tags) > 0 {
		tags := joinNonEmpty(n.Tags)
		if tags != "" {
			b.WriteString(" [tags: ")
			b.WriteString(tags)
			b.WriteString("]")
		}
	}

	// if kws := keywordNames(n.Keywords); kws != "" {
	// 	b.WriteString(" [kw: ")
	// 	b.WriteString(kws)
	// 	b.WriteString("]")
	// }

	if n.Locked {
		b.WriteString(" [locked]")
	}

	return b.String()
}

func safeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "<unnamed>"
	}
	return s
}

func joinNonEmpty(ss []string) string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return strings.Join(out, ", ")
}

func keywordNames(kws []*pb.Keyword) string {
	if len(kws) == 0 {
		return ""
	}
	out := make([]string, 0, len(kws))
	for _, k := range kws {
		if k == nil {
			continue
		}
		w := strings.TrimSpace(k.Keyword)
		if w != "" {
			out = append(out, w)
		}
	}
	return strings.Join(out, ", ")
}
