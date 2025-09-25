// Additional tests to improve server coverage
package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test addCompositeHandler with directory conflict checking
func TestAPI_AddCompositeHandler_ConflictChecking(t *testing.T) {
	// Setup: Create a temporary directory structure
	tmp := t.TempDir()
	parentDir := filepath.Join(tmp, "parent")
	childDir := filepath.Join(parentDir, "child")
	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test 1: Add initial manager successfully (assuming AddManager works)
	Composites = []*Folder{} // Reset composites

	// Instead of mocking AddManager, directly append to Composites for testing
	// Add initial manager directly
	// Composites = append(Composites, &Folder{Name: "parentManager", Path: parentDir})

	req := httptest.NewRequest("GET", "/addDirectory?name=parentManager&path="+parentDir, nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)

	if w.Body.String() != "true" {
		t.Errorf("Expected 'true' for successful addition, got %s", w.Body.String())
	}

	// Test 2: Try to add a child directory (should conflict)
	req = httptest.NewRequest("GET", "/addDirectory?name=childManager&path="+childDir, nil)
	w = httptest.NewRecorder()
	addCompositeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for conflict, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "already contained within existing manager") {
		t.Errorf("Expected conflict message, got: %s", responseBody)
	}

	// Test 3: Try to add duplicate manager name
	req = httptest.NewRequest("GET", "/addDirectory?name=parentManager&path="+tmp, nil)
	w = httptest.NewRecorder()
	addCompositeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate name, got %d", w.Code)
	}

	responseBody = w.Body.String()
	if !strings.Contains(responseBody, "name already exists") {
		t.Errorf("Expected duplicate name message, got: %s", responseBody)
	}
}

// Test the path containment helper functions
func TestAPI_PathContainmentHelpers(t *testing.T) {
	tests := []struct {
		parent   string
		child    string
		expected bool
		desc     string
	}{
		{"/home/user", "/home/user/documents", true, "child should be contained in parent"},
		{"/home/user/documents", "/home/user", false, "parent should not be contained in child"},
		{"/home/user", "/home/user", true, "identical paths should be contained"},
		{"/home/user", "/home/other", false, "unrelated paths should not be contained"},
		{"/home/user/docs", "/home/user/documents", false, "similar but different paths should not be contained"},
	}

	for _, tt := range tests {
		result := isPathContained(tt.parent, tt.child)
		if result != tt.expected {
			t.Errorf("%s: isPathContained(%q, %q) = %v, expected %v",
				tt.desc, tt.parent, tt.child, result, tt.expected)
		}
	}
}

// Test checkDirectoryConflicts function
func TestAPI_CheckDirectoryConflicts(t *testing.T) {
	tmp := t.TempDir()
	existingPath := filepath.Join(tmp, "existing")
	if err := os.MkdirAll(existingPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Setup existing composite
	Composites = []*Folder{{
		Name: "existing",
		Path: existingPath,
	}}

	// Test no conflict
	newPath := filepath.Join(tmp, "separate")
	hasConflict, message, err := checkDirectoryConflicts(newPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if hasConflict {
		t.Errorf("Expected no conflict for separate paths, got: %s", message)
	}

	// Test child conflict
	childPath := filepath.Join(existingPath, "child")
	hasConflict, message, err = checkDirectoryConflicts(childPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !hasConflict {
		t.Error("Expected conflict for child path")
	}
	if !strings.Contains(message, "already contained within") {
		t.Errorf("Expected containment message, got: %s", message)
	}

	// Test parent conflict
	parentPath := filepath.Dir(existingPath)
	hasConflict, message, err = checkDirectoryConflicts(parentPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !hasConflict {
		t.Error("Expected conflict for parent path")
	}
	if !strings.Contains(message, "would contain existing manager") {
		t.Errorf("Expected parent containment message, got: %s", message)
	}

	// Test exact match
	hasConflict, message, err = checkDirectoryConflicts(existingPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !hasConflict {
		t.Error("Expected conflict for exact path match")
	}
	if !strings.Contains(message, "already managed by") {
		t.Errorf("Expected exact match message, got: %s", message)
	}
}

// Test lock/unlock handlers edge cases
func TestAPI_LockUnlockHandlers_EdgeCases(t *testing.T) {
	// Setup test data
	tmp := t.TempDir()
	testFile := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	testFolder := &Folder{
		Name: "lockTest",
		Path: tmp,
		Files: []*File{{
			Name: "test.txt",
			Path: testFile,
		}},
	}
	Composites = []*Folder{testFolder}

	// Test missing parameters
	req := httptest.NewRequest("GET", "/lock?path="+tmp, nil) // missing name
	w := httptest.NewRecorder()
	lockHandler(w, req)
	if w.Body.String() != "Parameter missing" {
		t.Errorf("Expected 'Parameter missing', got %s", w.Body.String())
	}

	req = httptest.NewRequest("GET", "/lock?name=lockTest", nil) // missing path
	w = httptest.NewRecorder()
	lockHandler(w, req)
	if w.Body.String() != "Parameter missing" {
		t.Errorf("Expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test non-existent manager
	req = httptest.NewRequest("GET", "/lock?name=nonexistent&path="+tmp, nil)
	w = httptest.NewRecorder()
	lockHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("Expected 'false' for non-existent manager, got %s", w.Body.String())
	}

	// Test successful lock
	req = httptest.NewRequest("GET", "/lock?name=lockTest&path="+tmp, nil)
	w = httptest.NewRecorder()
	lockHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("Expected 'true' for successful lock, got %s", w.Body.String())
	}

	// Verify file was locked
	file := testFolder.GetFile(testFile)
	if file == nil || !file.Locked {
		t.Error("Expected file to be locked")
	}

	// Test unlock with same edge cases
	req = httptest.NewRequest("GET", "/unlock?path="+tmp, nil) // missing name
	w = httptest.NewRecorder()
	unlockHandler(w, req)
	if w.Body.String() != "Parameter missing" {
		t.Errorf("Expected 'Parameter missing', got %s", w.Body.String())
	}

	// Test successful unlock
	req = httptest.NewRequest("GET", "/unlock?name=lockTest&path="+tmp, nil)
	w = httptest.NewRecorder()
	unlockHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("Expected 'true' for successful unlock, got %s", w.Body.String())
	}

	// Verify file was unlocked
	if file.Locked {
		t.Error("Expected file to be unlocked")
	}
}

// Test addTagHandler edge cases
func TestAPI_AddTagHandler_EdgeCases(t *testing.T) {
	tmp := t.TempDir()
	testFile := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	testFolder := &Folder{
		Name: "tagTest",
		Path: tmp,
		Files: []*File{{
			Name: "test.txt",
			Path: testFile,
		}},
	}
	Composites = []*Folder{testFolder}

	// Test missing parameters
	req := httptest.NewRequest("GET", "/addTag?path="+testFile, nil) // missing tag
	w := httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("Expected 'false' for missing tag, got %s", w.Body.String())
	}

	req = httptest.NewRequest("GET", "/addTag?tag=testtag", nil) // missing path
	w = httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("Expected 'false' for missing path, got %s", w.Body.String())
	}

	// Test non-existent file
	req = httptest.NewRequest("GET", "/addTag?path=/nonexistent/file.txt&tag=testtag", nil)
	w = httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("Expected 'false' for non-existent file, got %s", w.Body.String())
	}

	// Test successful tag addition
	req = httptest.NewRequest("GET", "/addTag?path="+testFile+"&tag=testtag", nil)
	w = httptest.NewRecorder()
	addTagHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("Expected 'true' for successful tag addition, got %s", w.Body.String())
	}

	// Verify tag was added
	file := testFolder.GetFile(testFile)
	if file == nil {
		t.Fatal("File not found")
	}
	found := false
	for _, tag := range file.Tags {
		if tag == "testtag" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected tag to be added to file")
	}
}

// Test removeTagHandler with folder tags
func TestAPI_RemoveTagHandler_FolderTags(t *testing.T) {
	tmp := t.TempDir()
	subDir := filepath.Join(tmp, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	subfolder := &Folder{
		Name: "subdir",
		Path: subDir,
		Tags: []string{"foldertag"},
	}

	testFolder := &Folder{
		Name:       "tagTest",
		Path:       tmp,
		Subfolders: []*Folder{subfolder},
	}
	Composites = []*Folder{testFolder}

	// Test successful folder tag removal
	req := httptest.NewRequest("GET", "/removeTag?path="+subDir+"&tag=foldertag", nil)
	w := httptest.NewRecorder()
	removeTagHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("Expected 'true' for successful folder tag removal, got %s", w.Body.String())
	}

	// Verify tag was removed
	if len(subfolder.Tags) != 0 {
		t.Error("Expected folder tag to be removed")
	}

	// Test removing non-existent tag
	req = httptest.NewRequest("GET", "/removeTag?path="+subDir+"&tag=nonexistenttag", nil)
	w = httptest.NewRecorder()
	removeTagHandler(w, req)
	if w.Body.String() != "false" {
		t.Errorf("Expected 'false' for non-existent tag removal, got %s", w.Body.String())
	}
}

// Test deleteManager edge cases
func TestAPI_DeleteManagerHandler_EdgeCases(t *testing.T) {
	// Test with empty composites
	Composites = []*Folder{}

	req := httptest.NewRequest("GET", "/deleteManager?name=nonexistent", nil)
	w := httptest.NewRecorder()
	deleteManagerHandler(w, req)
	if w.Body.String() != "Manager not found" {
		t.Errorf("Expected 'Manager not found', got %s", w.Body.String())
	}

	// Test missing name parameter
	req = httptest.NewRequest("GET", "/deleteManager", nil)
	w = httptest.NewRecorder()
	deleteManagerHandler(w, req)
	if w.Body.String() != "Parameter missing" {
		t.Errorf("Expected 'Parameter missing', got %s", w.Body.String())
	}
}

// Test GetComposites function
func TestAPI_GetComposites(t *testing.T) {
	// Setup test data
	testComposite := &Folder{Name: "test", Path: "/test"}
	Composites = []*Folder{testComposite}

	// Test getter function
	result := GetComposites()
	if len(result) != 1 {
		t.Errorf("Expected 1 composite, got %d", len(result))
	}
	if result[0].Name != "test" {
		t.Errorf("Expected name 'test', got %s", result[0].Name)
	}
}

// Test concurrent access to handlers (mutex testing)
func TestAPI_ConcurrentAccess(t *testing.T) {
	tmp := t.TempDir()
	testFile := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	testFolder := &Folder{
		Name: "concurrentTest",
		Path: tmp,
		Files: []*File{{
			Name: "test.txt",
			Path: testFile,
		}},
	}
	Composites = []*Folder{testFolder}

	// Launch multiple goroutines to test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			tagName := fmt.Sprintf("tag%d", i)
			req := httptest.NewRequest("GET", "/addTag?path="+testFile+"&tag="+tagName, nil)
			w := httptest.NewRecorder()
			addTagHandler(w, req)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that some tags were added (exact number may vary due to race conditions)
	file := testFolder.GetFile(testFile)
	if file == nil {
		t.Fatal("File not found")
	}
	if len(file.Tags) == 0 {
		t.Error("Expected at least some tags to be added")
	}
}

// Test error handling in JSON encoding (edge case)
func TestAPI_JSONEncodingError(t *testing.T) {
	// This is harder to test directly, but we can at least ensure
	// the handlers set proper content types and handle basic cases

	tmp := t.TempDir()
	testFolder := &Folder{
		Name: "jsonTest",
		Path: tmp,
	}
	Composites = []*Folder{testFolder}

	// Test delete handlers that return JSON
	req := httptest.NewRequest("GET", "/deleteFile?name=jsonTest&path="+filepath.Join(tmp, "nonexistent.txt"), nil)
	w := httptest.NewRecorder()
	deleteFileHandler(w, req)

	// The handler should attempt to encode JSON even on error
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
}
func TestAPI_DeleteManagerHandler_Success(t *testing.T) {
	tmp := t.TempDir()
	testFolder := &Folder{Name: "delTest", Path: tmp}
	Composites = []*Folder{testFolder}
	managersFilePath = filepath.Join(tmp, "main.json")
	// Write a valid JSON file with the manager record
	recs := []ManagerRecord{{Name: "delTest", Path: tmp}}
	data, _ := json.Marshal(recs)
	os.WriteFile(managersFilePath, data, 0644)

	req := httptest.NewRequest("GET", "/deleteManager?name=delTest", nil)
	w := httptest.NewRecorder()
	deleteManagerHandler(w, req)
	if w.Body.String() != "true" {
		t.Errorf("Expected 'true', got %s", w.Body.String())
	}
	if len(Composites) != 0 {
		t.Errorf("Expected composites to be empty after deletion")
	}
}

func TestAPI_DeleteFolderHandler_Success(t *testing.T) {
	tmp := t.TempDir()
	subDir := filepath.Join(tmp, "subdir")
	os.MkdirAll(subDir, 0755)
	testFolder := &Folder{Name: "delFolderTest", Path: tmp, Subfolders: []*Folder{{Name: "subdir", Path: subDir}}}
	Composites = []*Folder{testFolder}

	req := httptest.NewRequest("GET", "/deleteFolder?name=delFolderTest&path="+subDir, nil)
	w := httptest.NewRecorder()
	deleteFolderHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if _, err := os.Stat(subDir); !os.IsNotExist(err) {
		t.Errorf("Expected subdir to be deleted")
	}
}
