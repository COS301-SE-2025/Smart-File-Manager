package filesystem

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAPI_AddAndRemoveDirectory(t *testing.T) {
	req := httptest.NewRequest("GET", "/addDirectory?name=testdir&path=../../testRootFolder", nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("expected true, got %s", w.Body.String())
	}

	req = httptest.NewRequest("GET", "/removeDirectory?path=../../testRootFolder", nil)
	w = httptest.NewRecorder()
	removeCompositeHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("expected true, got %s", w.Body.String())
	}
}

func TestAPI_AddAndRemoveTag(t *testing.T) {
	// Add folder to operate on
	req := httptest.NewRequest("GET", "/addDirectory?name=testdir&path=../../testRootFolder", nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)

	// Tag a file inside the test folder
	req = httptest.NewRequest("GET", "/addTag?path=../../testRootFolder/subdir/rb24.rs&tag=important", nil)
	w = httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("addTagHandler: expected true, got %s", w.Body.String())
	}

	// Remove the tag
	req = httptest.NewRequest("GET", "/removeTag?path=../../testRootFolder/subdir/rb24.rs&tag=important", nil)
	w = httptest.NewRecorder()
	removeTagHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("removeTagHandler: expected true, got %s", w.Body.String())
	}
}

func TestAPI_LockUnlock(t *testing.T) {
	// Add composite
	req := httptest.NewRequest("GET", "/addDirectory?name=testlock&path=../../testRootFolder", nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)

	// Lock
	req = httptest.NewRequest("GET", "/lock?path=../../testRootFolder", nil)
	w = httptest.NewRecorder()
	lockHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("lockHandler: expected true, got %s", w.Body.String())
	}

	// Unlock
	req = httptest.NewRequest("GET", "/unlock?path=../../testRootFolder", nil)
	w = httptest.NewRecorder()
	unlockHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("unlockHandler: expected true, got %s", w.Body.String())
	}

	// Cleanup
	req = httptest.NewRequest("GET", "/removeDirectory?path=../t../estRootFolder", nil)
	w = httptest.NewRecorder()
	removeCompositeHandler(w, req)
}

func TestAPI_EndpointsInvalidCases(t *testing.T) {
	// Try removing non-existent composite
	req := httptest.NewRequest("GET", "/removeDirectory?path=./invalid", nil)
	w := httptest.NewRecorder()
	removeCompositeHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("expected true even for non-existent remove, got %s", w.Body.String())
	}

	// Add tag to nonexistent file
	req = httptest.NewRequest("GET", "/addTag?path=./invalid/file.txt&tag=none", nil)
	w = httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("expected false for non-existent file tag, got %s", w.Body.String())
	}
}
func TestAPI_MoveDirectoryHandler(t *testing.T) {
	// 1) Create a temp project root called "Smart-File-Manager"
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		t.Fatal(err)
	}

	// 2) Inside it, make a simple source folder with one file
	srcDir := filepath.Join(projectRoot, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	fileName := "hello.txt"
	srcPath := filepath.Join(srcDir, fileName)
	content := []byte("content")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// 3) Prepare the composite list with exactly one Folder
	composites = []*Folder{{
		Name:    "testmgr",
		Path:    "src", // used by CreateDirectoryStructureRecursive
		NewPath: "src", // unused here
		Files: []*File{{
			Path:    filepath.Join("src", fileName),
			NewPath: fileName, // we want it at archives/testmgr/hello.txt
		}},
	}}

	// 4) chdir into projectRoot so getPath() will locate it
	origWd, _ := os.Getwd()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// 5) Ensure an empty archives folder exists
	if err := os.RemoveAll("archives"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("archives", 0755); err != nil {
		t.Fatal(err)
	}

	// 6) Call the handler
	req := httptest.NewRequest("GET", "/moveDirectory?name=testmgr", nil)
	w := httptest.NewRecorder()
	moveDirectoryHandler(w, req)

	// 7) Check HTTP response
	if got := w.Body.String(); got != "true" {
		t.Fatalf("moveDirectoryHandler returned %q; want \"true\"", got)
	}

	// 8) Verify the file was moved into archives/testmgr/hello.txt
	destPath := filepath.Join(projectRoot, "archives", "testmgr", fileName)
	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("moved file not found at %s: %v", destPath, err)
	}
	if string(data) != string(content) {
		t.Errorf("moved file content = %q; want %q", data, content)
	}
}
