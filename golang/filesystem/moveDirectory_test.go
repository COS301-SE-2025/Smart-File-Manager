package filesystem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreateDirectoryStructure(t *testing.T) {
	// Step 1: Clean up any existing archives directory before test
	err := os.RemoveAll("archives")
	if err != nil {
		t.Fatalf("failed to remove existing archives directory: %v", err)
	}

	err = os.Mkdir("archives", 0755)
	if err != nil {
		t.Fatalf("failed to create fresh archives directory: %v", err)
	}

	// Step 2: Ensure archives directory is removed after test
	t.Cleanup(func() {
		_ = os.RemoveAll("archives")
	})

	// Step 3: Create a mock folder structure
	managers := []string{"manager1", "manager2", "manager3"}
	var folders []*Folder
	for _, manager := range managers {
		folder := mockFolderStructureNamed(manager)
		folders = append(folders, folder)
	}
	for _, folder := range folders {
		CreateDirectoryStructure(folder)
	}

	// 4) walk the generated tree under actual project archives
	projectRoot := findProjectRoot(t)
	archives := filepath.Join(projectRoot, "archives")

	var got []string
	err = filepath.Walk(archives, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			rel, err := filepath.Rel(archives, path)
			if err != nil {
				return err
			}
			got = append(got, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking archive directory: %v", err)
	}

	// Step 5: Define the expected structure
	want := []string{
		".",
		"manager1",
		"manager1/test_root",
		"manager1/test_root/sub1",
		"manager1/test_root/sub1/sub1_1",
		"manager1/test_root/sub2",
		"manager2",
		"manager2/test_root",
		"manager2/test_root/sub1",
		"manager2/test_root/sub1/sub1_1",
		"manager2/test_root/sub2",
		"manager3",
		"manager3/test_root",
		"manager3/test_root/sub1",
		"manager3/test_root/sub1/sub1_1",
		"manager3/test_root/sub2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("directory tree mismatch:\n got:  %#v\n want: %#v", got, want)
	}
}

func TestMoveContent(t *testing.T) {
	projectRoot := findProjectRoot(t)
	archivesDir := filepath.Join(projectRoot, "archives")
	filesysDir := filepath.Join(projectRoot, "golang", "filesystem")

	// Clean up old test state
	_ = os.RemoveAll(archivesDir)
	_ = os.MkdirAll(filepath.Join(archivesDir, "myManager", "greeting"), 0755)

	// Create a dummy file in filesystem/ directory
	srcFilename := "hello.txt"
	srcPath := filepath.Join(filesysDir, srcFilename)
	content := []byte("👋 world")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcPath)

	// Folder to move
	item := &Folder{
		Name: "myManager",
		Files: []*File{
			{
				Name:    srcFilename,
				Path:    srcFilename, // relative to filesysDir
				NewPath: "greeting/hi.txt",
			},
		},
	}

	// Change working directory to where `moveContent` expects to be (filesystem/)
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(filesysDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Manually set root so it doesn't call getPath
	root = filepath.Join(projectRoot, "archives", item.Name)

	// Run function under test
	moveContent(item)

	// Assert source no longer exists
	if _, err := os.Stat(srcFilename); !os.IsNotExist(err) {
		t.Errorf("expected source %s to be gone, got err=%v", srcFilename, err)
	}

	// Assert target exists and has correct content
	destPath := filepath.Join(root, item.Files[0].NewPath)
	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("dest content = %q; want %q", data, content)
	}
}

// helper functions
func mockFolderStructureNamed(managerName string) *Folder {
	return &Folder{
		Name:    managerName,
		NewPath: "test_root",
		Subfolders: []*Folder{
			{
				NewPath: "test_root/sub1",
				Subfolders: []*Folder{
					{NewPath: "test_root/sub1/sub1_1"},
				},
			},
			{NewPath: "test_root/sub2"},
		},
	}
}

func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if filepath.Base(dir) == "Smart-File-Manager" {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root")
		}
		dir = parent
	}
}
