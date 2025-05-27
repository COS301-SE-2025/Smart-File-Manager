package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExploreExistingDirectory(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	rootPath := "../../testRootFolder"

	// Check that the test directory exists before continuing
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		t.Fatalf("Directory %s does not exist. Please create it before running the test.", rootPath)
	}

	// Convert root directory to composite
	root := ConvertToComposite("001", "TestRoot", rootPath)

	// Simple validation
	if len(root.ContainedItems) == 0 {
		t.Error("Expected at least one item in the root folder, but got none.")
	}

	// Check for a specific expected structure
	expectedFiles := map[string]bool{
		"a.txt":         false,
		"subdir":        false,
		"empty.txt":     false,
		"metadata.webp": false,
	}

	for _, item := range root.ContainedItems {
		if item.GetPath() == filepath.Join(rootPath, "a.txt") {
			expectedFiles["a.txt"] = true
		}
		if folder, ok := item.(*Folder); ok && folder.ItemName == "subdir" {
			expectedFiles["subdir"] = true
			for _, subItem := range folder.ContainedItems {
				if subItem.GetPath() == filepath.Join(rootPath, "subdir", "empty.txt") {
					expectedFiles["empty.txt"] = true
				}
				if subItem.GetPath() == filepath.Join(rootPath, "subdir", "metadata.webp") {
					expectedFiles["metadata.webp"] = true
				}
			}
		}
	}
	fmt.Println("==============")
	root.Display(0)
	fmt.Println("==============")

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("Expected to find %s but did not", file)
		}
	}
	deleteComposite(&root)
}
