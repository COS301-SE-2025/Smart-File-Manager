package filesystem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Test simple duplicates in a single folder
func TestFindDuplicateFiles_Simple(t *testing.T) {
	tmp := t.TempDir()

	// Create two files with identical content
	pathA := filepath.Join(tmp, "fileA.txt")
	pathB := filepath.Join(tmp, "fileB.txt")
	content := []byte("duplicate content")
	if err := os.WriteFile(pathA, content, 0644); err != nil {
		t.Fatalf("failed to write file A: %v", err)
	}
	if err := os.WriteFile(pathB, content, 0644); err != nil {
		t.Fatalf("failed to write file B: %v", err)
	}

	// Build Folder struct
	folder := &Folder{
		Files: []*File{
			{Name: "fileA.txt", Path: pathA},
			{Name: "fileB.txt", Path: pathB},
		},
	}

	// Invoke duplicate finder
	dups := FindDuplicateFiles(folder)

	// Expect exactly one duplicate pair
	if len(dups) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(dups))
	}

	exp := DuplicateEntry{
		Name:      "fileB.txt",
		Original:  pathA,
		Duplicate: pathB,
	}
	if !reflect.DeepEqual(dups[0], exp) {
		t.Errorf("unexpected duplicate entry: got  %+v want %+v", dups[0], exp)
	}
}

// Test duplicates within nested subfolders
func TestFindDuplicateFiles_Nested(t *testing.T) {
	tmp := t.TempDir()

	// Create a subfolder and two duplicate files inside
	sub := filepath.Join(tmp, "subfolder")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatalf("failed to create subfolder: %v", err)
	}
	path1 := filepath.Join(sub, "a.txt")
	path2 := filepath.Join(sub, "copy_of_a.txt")
	content := []byte("nested duplicate")
	if err := os.WriteFile(path1, content, 0644); err != nil {
		t.Fatalf("failed to write path1: %v", err)
	}
	if err := os.WriteFile(path2, content, 0644); err != nil {
		t.Fatalf("failed to write path2: %v", err)
	}

	// Build Folder struct with nested Folder
	subFolder := &Folder{
		Files: []*File{
			{Name: "a.txt", Path: path1},
			{Name: "copy_of_a.txt", Path: path2},
		},
	}
	root := &Folder{
		Subfolders: []*Folder{subFolder},
	}

	// Invoke duplicate finder
	dups := FindDuplicateFiles(root)

	// Expect exactly one duplicate pair
	if len(dups) != 1 {
		t.Fatalf("expected 1 duplicate in nested, got %d", len(dups))
	}

	exp := DuplicateEntry{
		Name:      "copy_of_a.txt",
		Original:  path1,
		Duplicate: path2,
	}
	if !reflect.DeepEqual(dups[0], exp) {
		t.Errorf("unexpected nested duplicate entry: got  %+v want %+v", dups[0], exp)
	}
}
