package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConvertToObject(t *testing.T) {
	rootPath := "../../testRootFolder"
	subdirPath := filepath.Join(rootPath, "subdir")
	hiddenPath := filepath.Join(subdirPath, ".hiddenFolder")

	// Ensure testRootFolder/subdir exists
	if info, err := os.Stat(subdirPath); err != nil || !info.IsDir() {
		t.Fatalf("Expected subdir at %s, got error: %v", subdirPath, err)
	}

	// Create .hiddenFolder inside subdir for testing
	err := os.MkdirAll(hiddenPath, 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create hidden folder for test: %v", err)
	}
	defer os.RemoveAll(hiddenPath) // Cleanup after test

	// Convert directory into Folder tree
	root, err := ConvertToObject("TestRoot", rootPath)
	if err != nil {
		t.Fatalf("ConvertToObject returned error: %v", err)
	}

	// Expect at least one file or subfolder under root
	if len(root.Files)+len(root.Subfolders) == 0 {
		t.Error("Expected at least one file or subfolder in root, got none")
	}

	// Check top-level expected items
	expected := map[string]struct{}{
		"a.txt":  {},
		"subdir": {},
	}
	for _, f := range root.Files {
		if _, ok := expected[f.Name]; ok {
			delete(expected, f.Name)
		}
	}

	for _, sf := range root.Subfolders {
		if sf.Name == "subdir" {
			delete(expected, "subdir")

			// Check nested files
			nestedExpected := map[string]struct{}{
				"empty.txt":     {},
				"metadata.webp": {},
			}
			for _, nf := range sf.Files {
				if _, ok := nestedExpected[nf.Name]; ok {
					delete(nestedExpected, nf.Name)
				}
			}
			for name := range nestedExpected {
				t.Errorf("Expected to find %s in subdir, but did not", name)
			}

			// 🔒 Verify subdir is auto-locked due to hidden folder
			if !sf.Locked {
				t.Errorf("Expected folder %q to be locked due to hidden folder, but it was not", sf.Name)
			}
		}
	}

	// Report missing top-level files/folders
	for name := range expected {
		t.Errorf("Expected to find %s in root, but did not", name)
	}

	// Optionally visualize
	root.Display(0)
}
