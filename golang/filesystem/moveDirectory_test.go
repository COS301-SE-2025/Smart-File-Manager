package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

// Mock folder structure
func mockFolderStructure() *Folder {
	return &Folder{
		NewPath: "test_root",
		Subfolders: []*Folder{
			{
				NewPath: "test_root/sub1",
				Subfolders: []*Folder{
					{
						NewPath: "test_root/sub1/sub1_1",
					},
				},
			},
			{
				NewPath: "test_root/sub2",
			},
		},
	}
}

func TestCreateDirectoryStructure(t *testing.T) {
	// Create a temporary directory
	tmpRoot, err := os.MkdirTemp("", "sfm_test")
	if err != nil {
		t.Fatalf("Failed to create temp root directory: %v", err)
	}
	defer os.RemoveAll(tmpRoot) // Clean up afterwards

	// Override global root for testing
	root = tmpRoot

	// Prepare test data
	mockFolder := mockFolderStructure()

	// Run the function
	CreateDirectoryStructure(mockFolder)

	// Check that all expected directories exist
	expectedPaths := []string{
		filepath.Join(tmpRoot, "archives", "test_root"),
		filepath.Join(tmpRoot, "archives", "test_root/sub1"),
		filepath.Join(tmpRoot, "archives", "test_root/sub1/sub1_1"),
		filepath.Join(tmpRoot, "archives", "test_root/sub2"),
	}

	for _, path := range expectedPaths {
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Expected directory %s was not created: %v", path, err)
		} else if !info.IsDir() {
			t.Errorf("Expected path %s is not a directory", path)
		}
	}
}
