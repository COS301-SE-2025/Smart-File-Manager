package filesystem

import (
	"os"
	"testing"
)

func TestConvertToObject(t *testing.T) {
	// Define path to test root folder (adjust relative path as needed)
	rootPath := "../../testRootFolder"

	// Ensure test directory exists
	if info, err := os.Stat(rootPath); err != nil || !info.IsDir() {
		t.Fatalf("Test directory %s does not exist or is not a directory", rootPath)
	}

	// Convert directory into Folder tree
	root, err := ConvertToObject("TestRoot", rootPath)
	if err != nil {
		t.Fatalf("ConvertToObject returned error: %v", err)
	}

	// Expect at least one file or subfolder under root
	if len(root.Files)+len(root.Subfolders) == 0 {
		t.Error("Expected at least one file or subfolder in root, got none")
	}

	// Map of expected items
	expected := map[string]struct{}{
		"a.txt":  {},
		"subdir": {},
	}

	// Check root files
	for _, f := range root.Files {
		if _, ok := expected[f.Name]; ok {
			delete(expected, f.Name)
		}
	}

	// Find 'subdir' folder and inspect its contents
	for _, sf := range root.Subfolders {
		if sf.Name == "subdir" {
			// mark found
			delete(expected, "subdir")
			// look for nested items
			nestedExpected := map[string]struct{}{
				"empty.txt":     {},
				"metadata.webp": {},
			}
			for _, nf := range sf.Files {
				if _, ok := nestedExpected[nf.Name]; ok {
					delete(nestedExpected, nf.Name)
				}
			}
			// report missing in subdir
			for name := range nestedExpected {
				t.Errorf("Expected to find %s in subdir, but did not", name)
			}
		}
	}

	// Report missing in root
	for name := range expected {
		t.Errorf("Expected to find %s in root, but did not", name)
	}
	root.Display(0)
}
