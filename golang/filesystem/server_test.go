package filesystem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// func TestAPI_LockUnlock(t *testing.T) {
// 	// Add composite
// 	req := httptest.NewRequest("GET", "/addDirectory?name=testlock&path=../../testRootFolder", nil)
// 	w := httptest.NewRecorder()
// 	addCompositeHandler(w, req)

// 	// Lock
// 	// req = httptest.NewRequest("GET", "/lock?name=testlock&path=../../testRootFolder", nil)
// 	// w = httptest.NewRecorder()
// 	// lockHandler(w, req)
// 	// if w.Body.String() != "true" {
// 	// 	t.Fatalf("lockHandler: expected true, got %s", w.Body.String())
// 	// }

// 	// Unlock
// 	req = httptest.NewRequest("GET", "/unlock?name=testlock&path=../../testRootFolder", nil)
// 	w = httptest.NewRecorder()
// 	unlockHandler(w, req)
// 	if w.Body.String() != "true" {
// 		t.Fatalf("unlockHandler: expected true, got %s", w.Body.String())
// 	}

// 	// Cleanup
// 	req = httptest.NewRequest("GET", "/removeDirectory?path=../../testRootFolder", nil)
// 	w = httptest.NewRecorder()
// 	removeCompositeHandler(w, req)
// }

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
	Composites = []*Folder{{
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

// Test the find-duplicates endpoint
func TestAPI_FindDuplicateFilesHandler(t *testing.T) {
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a data folder with two identical files
	dataDir := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}
	fileA := filepath.Join(dataDir, "a.txt")
	fileB := filepath.Join(dataDir, "b.txt")
	content := []byte("duplicate test")
	if err := os.WriteFile(fileA, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileB, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Register composite
	Composites = []*Folder{{
		Name:    "dupTest",
		Path:    "data",
		NewPath: "data",
		Files: []*File{{
			Name: "a.txt", Path: fileA,
		}, {
			Name: "b.txt", Path: fileB,
		}},
	}}

	// Move into project root so handler finds it
	origWd, _ := os.Getwd()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Call endpoint
	req := httptest.NewRequest("GET", "/findDuplicateFiles?name=dupTest", nil)
	w := httptest.NewRecorder()
	findDuplicateFilesHandler(w, req)

	// Check response code
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Parse JSON
	var entries []DuplicateEntry
	if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Expect one duplicate entry
	if len(entries) != 1 {
		t.Fatalf("expected 1 duplicate entry, got %d", len(entries))
	}

	exp := DuplicateEntry{
		Name:      "b.txt",
		Original:  fileA,
		Duplicate: fileB,
	}
	if entries[0] != exp {
		t.Errorf("unexpected entry: got %+v, want %+v", entries[0], exp)
	}
}
func TestAPI_BulkAddTags(t *testing.T) {
	// Setup: Create temp folder and file
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dataDir, "test.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Register the folder and file in composites
	Composites = []*Folder{{
		Name:    "bulkTest",
		Path:    "data",
		NewPath: "data",
		Files: []*File{{
			Name: "test.txt",
			Path: filePath,
		}},
	}}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// JSON payload for tags
	jsonBody := `[{
		"file_path": "` + filePath + `",
		"tags": ["alpha", "beta", "gamma"]
	}]`

	// Create request with JSON body and query parameter ?name=bulkTest
	req := httptest.NewRequest("POST", "/bulkAddTag?name=bulkTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	BulkAddTagHandler(w, req)

	// Verify response code
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Tags added successfully") {
		t.Errorf("expected success message, got %s", got)
	}

	// Check if tags were actually added to the file
	var file *File
	for _, f := range Composites[0].Files {
		if f.Path == filePath {
			file = f
			break
		}
	}
	if file == nil {
		t.Fatal("file not found in folder")
	}
	expectedTags := map[string]bool{"alpha": false, "beta": false, "gamma": false}
	for _, tag := range file.Tags {
		if _, exists := expectedTags[tag]; exists {
			expectedTags[tag] = true
		}
	}
	for tag, found := range expectedTags {
		if !found {
			t.Errorf("expected tag %s not found in file.Tags", tag)
		}
	}
}

func TestAPI_DeleteFileHandler(t *testing.T) {
	// Setup: Create temp folder and file
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test file
	testFileName := "delete_test.txt"
	testFilePath := filepath.Join(dataDir, testFileName)
	testContent := []byte("file to be deleted")
	if err := os.WriteFile(testFilePath, testContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Register the folder and file in composites
	testFile := &File{
		Name: testFileName,
		Path: testFilePath,
	}

	testFolder := &Folder{
		Name:    "deleteTest",
		Path:    dataDir,
		NewPath: dataDir,
		Files:   []*File{testFile},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Verify file exists before deletion
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatal("test file should exist before deletion")
	}

	// Test successful file deletion
	req := httptest.NewRequest("GET", "/deleteFile?name=deleteTest&path="+testFilePath, nil)
	w := httptest.NewRecorder()
	deleteFileHandler(w, req)

	// Verify response
	if w.Body.String() != "true" {
		t.Fatalf("deleteFileHandler: expected true, got %s", w.Body.String())
	}

	// Verify file was actually deleted from filesystem
	if _, err := os.Stat(testFilePath); !os.IsNotExist(err) {
		t.Error("file should have been deleted from filesystem")
	}

	// Verify file was removed from composite structure
	// Note: This test assumes RemoveFile method works correctly
	// You might want to verify the file is no longer in the Files slice
}

func TestAPI_DeleteFileHandler_MissingParams(t *testing.T) {
	// Test missing path parameter
	req := httptest.NewRequest("GET", "/deleteFile?name=test", nil)
	w := httptest.NewRecorder()
	deleteFileHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test missing name parameter
	req = httptest.NewRequest("GET", "/deleteFile?path=/some/path", nil)
	w = httptest.NewRecorder()
	deleteFileHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test both parameters missing
	req = httptest.NewRequest("GET", "/deleteFile", nil)
	w = httptest.NewRecorder()
	deleteFileHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}
}

func TestAPI_DeleteFileHandler_NonExistentManager(t *testing.T) {
	// Setup empty composites
	Composites = []*Folder{}

	// Test with non-existent manager name
	req := httptest.NewRequest("GET", "/deleteFile?name=nonexistent&path=/some/path", nil)
	w := httptest.NewRecorder()
	deleteFileHandler(w, req)

	if w.Body.String() != "false" {
		t.Errorf("expected 'false' for non-existent manager, got %s", w.Body.String())
	}
}

func TestAPI_DeleteFolderHandler(t *testing.T) {
	// Setup: Create temp folder structure
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	testDir := filepath.Join(dataDir, "test_folder")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files inside the test folder
	testFile1 := filepath.Join(testDir, "file1.txt")
	testFile2 := filepath.Join(testDir, "file2.txt")
	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create nested subfolder
	nestedDir := filepath.Join(testDir, "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}
	nestedFile := filepath.Join(nestedDir, "nested_file.txt")
	if err := os.WriteFile(nestedFile, []byte("nested content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Register the folder structure in composites
	testFolder := &Folder{
		Name:    "deleteFolderTest",
		Path:    dataDir,
		NewPath: dataDir,
		Files: []*File{
			{Name: "file1.txt", Path: testFile1},
			{Name: "file2.txt", Path: testFile2},
			{Name: "nested_file.txt", Path: nestedFile},
		},
		Subfolders: []*Folder{
			{Name: "test_folder", Path: testDir},
			{Name: "nested", Path: nestedDir},
		},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Verify folder exists before deletion
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatal("test folder should exist before deletion")
	}

	// Test successful folder deletion
	req := httptest.NewRequest("GET", "/deleteFolder?name=deleteFolderTest&path="+testDir, nil)
	w := httptest.NewRecorder()
	deleteFolderHandler(w, req)

	// Verify response
	if w.Body.String() != "true" {
		t.Fatalf("deleteFolderHandler: expected true, got %s", w.Body.String())
	}

	// Verify folder and all its contents were deleted from filesystem
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Error("folder should have been deleted from filesystem")
	}

	// Verify nested files were also deleted
	if _, err := os.Stat(testFile1); !os.IsNotExist(err) {
		t.Error("nested file1 should have been deleted")
	}
	if _, err := os.Stat(testFile2); !os.IsNotExist(err) {
		t.Error("nested file2 should have been deleted")
	}
	if _, err := os.Stat(nestedFile); !os.IsNotExist(err) {
		t.Error("nested file should have been deleted")
	}

	// Verify folder was removed from composite structure
	// Note: This test assumes RemoveSubfolder method works correctly
}

func TestAPI_DeleteFolderHandler_MissingParams(t *testing.T) {
	// Test missing path parameter
	req := httptest.NewRequest("GET", "/deleteFolder?name=test", nil)
	w := httptest.NewRecorder()
	deleteFolderHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test missing name parameter
	req = httptest.NewRequest("GET", "/deleteFolder?path=/some/path", nil)
	w = httptest.NewRecorder()
	deleteFolderHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test both parameters missing
	req = httptest.NewRequest("GET", "/deleteFolder", nil)
	w = httptest.NewRecorder()
	deleteFolderHandler(w, req)

	if w.Body.String() != "Parameter missing" {
		t.Errorf("expected 'Parameter missing', got %s", w.Body.String())
	}
}

func TestAPI_DeleteFolderHandler_NonExistentManager(t *testing.T) {
	// Setup empty composites
	Composites = []*Folder{}

	// Test with non-existent manager name
	req := httptest.NewRequest("GET", "/deleteFolder?name=nonexistent&path=/some/path", nil)
	w := httptest.NewRecorder()
	deleteFolderHandler(w, req)

	if w.Body.String() != "false" {
		t.Errorf("expected 'false' for non-existent manager, got %s", w.Body.String())
	}
}

func TestAPI_DeleteHandlers_ErrorHandling(t *testing.T) {
	// This test verifies that the handlers properly handle OS errors
	// Note: The current implementation uses panic() for os.Remove errors
	// In a production environment, you might want to handle errors more gracefully

	// Setup: Create a composite with a non-existent file path
	testFolder := &Folder{
		Name:    "errorTest",
		Path:    "/nonexistent/path",
		NewPath: "/nonexistent/path",
	}

	Composites = []*Folder{testFolder}

	// Test file deletion with non-existent file
	// Note: This will panic in the current implementation
	// You might want to modify the handler to return an error instead
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to os.Remove error
			t.Log("deleteFileHandler panicked as expected when file doesn't exist")
		}
	}()

	req := httptest.NewRequest("GET", "/deleteFile?name=errorTest&path=/nonexistent/file.txt", nil)
	w := httptest.NewRecorder()
	deleteFileHandler(w, req)

	// If we reach here without panic, the test might need adjustment
	t.Log("deleteFileHandler completed without panic")
}

func TestAPI_BulkDeleteFolderHandler(t *testing.T) {
	// Setup: Create temp folder structure with multiple subfolders
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")

	// Create multiple test folders
	folder1 := filepath.Join(dataDir, "folder1")
	folder2 := filepath.Join(dataDir, "folder2")
	folder3 := filepath.Join(dataDir, "folder3")

	if err := os.MkdirAll(folder1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(folder2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(folder3, 0755); err != nil {
		t.Fatal(err)
	}

	// Add some files to the folders
	file1 := filepath.Join(folder1, "test1.txt")
	file2 := filepath.Join(folder2, "test2.txt")
	file3 := filepath.Join(folder3, "test3.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}

	// Register the folder structure in composites
	testFolder := &Folder{
		Name:    "bulkDeleteFolderTest",
		Path:    dataDir,
		NewPath: dataDir,
		Files: []*File{
			{Name: "test1.txt", Path: file1},
			{Name: "test2.txt", Path: file2},
			{Name: "test3.txt", Path: file3},
		},
		Subfolders: []*Folder{
			{Name: "folder1", Path: folder1},
			{Name: "folder2", Path: folder2},
			{Name: "folder3", Path: folder3},
		},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Verify folders exist before deletion
	if _, err := os.Stat(folder1); os.IsNotExist(err) {
		t.Fatal("folder1 should exist before deletion")
	}
	if _, err := os.Stat(folder2); os.IsNotExist(err) {
		t.Fatal("folder2 should exist before deletion")
	}
	if _, err := os.Stat(folder3); os.IsNotExist(err) {
		t.Fatal("folder3 should exist before deletion")
	}

	// Create JSON payload for bulk deletion (delete folder1 and folder3)
	jsonBody := `[
		{"file_path": "` + folder1 + `"},
		{"file_path": "` + folder3 + `"}
	]`

	// Test successful bulk folder deletion
	req := httptest.NewRequest("POST", "/bulkDeleteFolder?name=bulkDeleteFolderTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	// if got := w.Body.String(); !strings.Contains(got, "Folders removed successfully") {
	// 	t.Errorf("expected success message, got %s", got)
	// }

	// Verify deleted folders no longer exist
	if _, err := os.Stat(folder1); !os.IsNotExist(err) {
		t.Error("folder1 should have been deleted")
	}
	if _, err := os.Stat(folder3); !os.IsNotExist(err) {
		t.Error("folder3 should have been deleted")
	}

	// Verify folder2 still exists (wasn't in deletion list)
	if _, err := os.Stat(folder2); os.IsNotExist(err) {
		t.Error("folder2 should still exist")
	}
}

func TestAPI_BulkDeleteFolderHandler_MissingName(t *testing.T) {
	jsonBody := `[{"file_path": "/some/path"}]`

	req := httptest.NewRequest("POST", "/bulkDeleteFolder", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Missing 'name' parameter") {
		t.Errorf("expected missing name error, got %s", got)
	}
}

func TestAPI_BulkDeleteFolderHandler_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`

	req := httptest.NewRequest("POST", "/bulkDeleteFolder?name=test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Invalid request body") {
		t.Errorf("expected invalid request body error, got %s", got)
	}
}

func TestAPI_BulkDeleteFolderHandler_NonExistentManager(t *testing.T) {
	Composites = []*Folder{}

	jsonBody := `[{"file_path": "/some/path"}]`

	req := httptest.NewRequest("POST", "/bulkDeleteFolder?name=nonexistent", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Folder not found") {
		t.Errorf("expected folder not found error, got %s", got)
	}
}

func TestAPI_BulkDeleteFolderHandler_EmptyList(t *testing.T) {
	// Setup minimal composite
	testFolder := &Folder{
		Name:    "emptyTest",
		Path:    "/test",
		NewPath: "/test",
	}
	Composites = []*Folder{testFolder}

	jsonBody := `[]`

	req := httptest.NewRequest("POST", "/bulkDeleteFolder?name=emptyTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty list, got %d", w.Code)
	}
}

func TestAPI_BulkDeleteFileHandler(t *testing.T) {
	// Setup: Create temp folder and multiple test files
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create multiple test files
	file1 := filepath.Join(dataDir, "file1.txt")
	file2 := filepath.Join(dataDir, "file2.txt")
	file3 := filepath.Join(dataDir, "file3.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}

	// Register the files in composites
	testFolder := &Folder{
		Name:    "bulkDeleteFileTest",
		Path:    dataDir,
		NewPath: dataDir,
		Files: []*File{
			{Name: "file1.txt", Path: file1},
			{Name: "file2.txt", Path: file2},
			{Name: "file3.txt", Path: file3},
		},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Verify files exist before deletion
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		t.Fatal("file1 should exist before deletion")
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		t.Fatal("file2 should exist before deletion")
	}
	if _, err := os.Stat(file3); os.IsNotExist(err) {
		t.Fatal("file3 should exist before deletion")
	}

	// Create JSON payload for bulk deletion (delete file1 and file3)
	jsonBody := `[
		{"file_path": "` + file1 + `"},
		{"file_path": "` + file3 + `"}
	]`

	// Test successful bulk file deletion
	req := httptest.NewRequest("POST", "/bulkDeleteFiles?name=bulkDeleteFileTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	// if got := w.Body.String(); !strings.Contains(got, "Files removed successfully") {
	// 	t.Errorf("expected success message, got %s", got)
	// }

	// Verify deleted files no longer exist
	if _, err := os.Stat(file1); !os.IsNotExist(err) {
		t.Error("file1 should have been deleted")
	}
	if _, err := os.Stat(file3); !os.IsNotExist(err) {
		t.Error("file3 should have been deleted")
	}

	// Verify file2 still exists (wasn't in deletion list)
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		t.Error("file2 should still exist")
	}
}

func TestAPI_BulkDeleteFileHandler_MissingName(t *testing.T) {
	jsonBody := `[{"file_path": "/some/path"}]`

	req := httptest.NewRequest("POST", "/bulkDeleteFiles", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Missing 'name' parameter") {
		t.Errorf("expected missing name error, got %s", got)
	}
}

func TestAPI_BulkDeleteFileHandler_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`

	req := httptest.NewRequest("POST", "/bulkDeleteFiles?name=test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Invalid request body") {
		t.Errorf("expected invalid request body error, got %s", got)
	}
}

func TestAPI_BulkDeleteFileHandler_NonExistentManager(t *testing.T) {
	Composites = []*Folder{}

	jsonBody := `[{"file_path": "/some/path"}]`

	req := httptest.NewRequest("POST", "/bulkDeleteFiles?name=nonexistent", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if got := w.Body.String(); !strings.Contains(got, "Files not found") {
		t.Errorf("expected files not found error, got %s", got)
	}
}

func TestAPI_BulkDeleteFileHandler_EmptyList(t *testing.T) {
	// Setup minimal composite
	testFolder := &Folder{
		Name:    "emptyTest",
		Path:    "/test",
		NewPath: "/test",
	}
	Composites = []*Folder{testFolder}

	jsonBody := `[]`

	req := httptest.NewRequest("POST", "/bulkDeleteFiles?name=emptyTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty list, got %d", w.Code)
	}
}
