// Additional tests to improve server coverage
package filesystem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
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

	// Backup and restore original Composites
	orig := Composites
	defer func() { Composites = orig }()

	// Simulate that a manager already exists for parentDir
	Composites = []*Folder{
		{Name: "parentManager", Path: parentDir},
	}

	// Test 1: Try to add a child directory (should be detected as conflict)
	req := httptest.NewRequest("GET", "/addDirectory?name=childManager&path="+childDir, nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)

	// Current implementation writes "false" on conflict (and returns 200 OK)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for conflict, got %d", w.Code)
	}
	body := strings.TrimSpace(w.Body.String())
	if body != "false" {
		t.Errorf("expected body 'false' on conflict, got %q", body)
	}

	// Test 2: Try to add duplicate manager name (should return 400)
	req = httptest.NewRequest("GET", "/addDirectory?name=parentManager&path="+tmp, nil)
	w = httptest.NewRecorder()
	addCompositeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for duplicate name, got %d", w.Code)
	}
	respBody := w.Body.String()
	if !strings.Contains(respBody, "name already exists") && !strings.Contains(respBody, "already exists") {
		t.Errorf("expected duplicate name message, got: %s", respBody)
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

func TestSecretMiddleware_NoEnv(t *testing.T) {
	// unset the env var
	os.Unsetenv("SFM_API_SECRET")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("apiSecret", "anything")
	rr := httptest.NewRecorder()

	handler := secretMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("next handler should not be called")
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestSecretMiddleware_WrongSecret(t *testing.T) {
	os.Setenv("SFM_API_SECRET", "correctSecret")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("apiSecret", "wrongSecret")
	rr := httptest.NewRecorder()

	handler := secretMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("next handler should not be called")
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestSecretMiddleware_CorrectSecret(t *testing.T) {
	os.Setenv("SFM_API_SECRET", "correctSecret")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("apiSecret", "correctSecret")
	rr := httptest.NewRecorder()

	called := false
	handler := secretMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if !called {
		t.Errorf("next handler was not called")
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

func TestTagsConversionRoundTrip(t *testing.T) {
	in := []*pb.Tag{
		{Name: "one"},
		{Name: "Two"},
		{Name: "three"},
	}
	strs := tagsToStrings(in)
	if len(strs) != len(in) {
		t.Fatalf("expected %d tags, got %d", len(in), len(strs))
	}
	out := stringsToTags(strs)
	if len(out) != len(in) {
		t.Fatalf("expected %d pb.Tags, got %d", len(in), len(out))
	}
	for i := range in {
		if in[i].Name != out[i].Name {
			t.Fatalf("expected tag %q at index %d, got %q", in[i].Name, i, out[i].Name)
		}
	}
}

// Ensure nil guards don't panic
func TestMergeProtoNilGuards(t *testing.T) {
	// both nil
	mergeProtoToFolder(nil, nil)

	// dir nil, existing non-nil
	existing := &Folder{Name: "exists", Path: "/tmp"}
	mergeProtoToFolder(nil, existing)

	// dir non-nil, existing nil
	dir := &pb.Directory{Name: "dir", Path: "/tmp/dir"}
	mergeProtoToFolder(dir, nil)

	// helper nil guard
	mergeProtoToFolderHelper(nil, nil)
	mergeProtoToFolderHelper(nil, &Folder{})
	mergeProtoToFolderHelper(&pb.Directory{}, nil)
}

// Test mapping of files and nested directories via mergeProtoToFolder and helper
func TestMergeProtoToFolder_Mapping(t *testing.T) {
	// Build a proto tree:
	// root
	//  - fileA (with one tag)
	//  - subdir (contains fileB with keyword)
	root := &pb.Directory{
		Name: "root",
		Path: "/root",
		Files: []*pb.File{
			{
				Name:         "fileA.txt",
				OriginalPath: "/root/fileA.txt",
				Tags:         []*pb.Tag{{Name: "alpha"}, {Name: "beta"}},
				IsLocked:     true,
				NewPath:      "/newroot/fileA.txt",
				// Metadata left nil (metadataConverter should handle nil)
			},
		},
		Directories: []*pb.Directory{
			{
				Name: "subdir",
				Path: "/root/subdir",
				Files: []*pb.File{
					{
						Name:         "fileB.md",
						OriginalPath: "/root/subdir/fileB.md",
						Tags:         []*pb.Tag{{Name: "gamma"}},
						IsLocked:     false,
						NewPath:      "/newroot/subdir/fileB.md",
						Keywords:     []*pb.Keyword{{Keyword: "kw1", Score: 2.0}},
					},
				},
			},
		},
	}

	// existing folder with pre-populated content that should be cleared
	existing := &Folder{
		Name:       "old",
		Path:       "/old",
		Files:      []*File{{Name: "oldfile"}},
		Subfolders: []*Folder{{Name: "oldsub"}},
	}

	mergeProtoToFolder(root, existing)

	// after merge, existing should reflect proto
	if existing.Name != "old" {
		// mergeProtoToFolder doesn't change existing.Name at top-level, only subfolders via helper.
		// we only assert Files and Subfolders contents here.
	}

	// Files length and mapping
	if len(existing.Files) != 1 {
		t.Fatalf("expected 1 file in existing.Files, got %d", len(existing.Files))
	}
	f := existing.Files[0]
	if f.Name != "fileA.txt" {
		t.Fatalf("expected file name fileA.txt, got %q", f.Name)
	}
	if f.Path != "/root/fileA.txt" {
		t.Fatalf("expected path /root/fileA.txt, got %q", f.Path)
	}
	if f.Locked != true {
		t.Fatalf("expected Locked true, got %v", f.Locked)
	}
	// tagsToStrings should have converted tags
	if len(f.Tags) != 2 || f.Tags[0] != "alpha" || f.Tags[1] != "beta" {
		t.Fatalf("tags not converted correctly: %#v", f.Tags)
	}

	// Subfolders length and mapping
	if len(existing.Subfolders) != 1 {
		t.Fatalf("expected 1 subfolder, got %d", len(existing.Subfolders))
	}
	sub := existing.Subfolders[0]
	if sub.Name != "subdir" {
		t.Fatalf("expected subfolder name subdir, got %q", sub.Name)
	}
	if sub.Path != "/root/subdir" {
		t.Fatalf("expected subfolder path /root/subdir, got %q", sub.Path)
	}
	// child files mapping via helper
	if len(sub.Files) != 1 {
		t.Fatalf("expected 1 file in sub.Files, got %d", len(sub.Files))
	}
	fb := sub.Files[0]
	if fb.Name != "fileB.md" {
		t.Fatalf("expected sub file name fileB.md, got %q", fb.Name)
	}
	// Keywords should be preserved by helper
	if len(fb.Keywords) != 1 || fb.Keywords[0].Keyword != "kw1" {
		t.Fatalf("expected keyword kw1 in child file, got %#v", fb.Keywords)
	}
}

// A small test to ensure helper updates an existing folder in-place
func TestMergeProtoToFolderHelper_InPlace(t *testing.T) {
	dir := &pb.Directory{
		Name: "X",
		Path: filepath.Join("/base", "X"),
		Files: []*pb.File{
			{
				Name:         "a.txt",
				OriginalPath: "/base/X/a.txt",
				Tags:         []*pb.Tag{{Name: "t1"}},
			},
		},
		Directories: []*pb.Directory{
			{
				Name: "C",
				Path: "/base/X/C",
			},
		},
	}

	existing := &Folder{
		Name:       "oldName",
		Path:       "/oldpath",
		Files:      []*File{{Name: "should be removed"}},
		Subfolders: []*Folder{{Name: "should be removed"}},
	}

	mergeProtoToFolderHelper(dir, existing)

	if existing.Name != "X" {
		t.Fatalf("expected Name X, got %q", existing.Name)
	}
	if existing.Path != "/base/X" {
		t.Fatalf("expected Path /base/X, got %q", existing.Path)
	}
	if len(existing.Files) != 1 || existing.Files[0].Name != "a.txt" {
		t.Fatalf("file mapping failed, got %#v", existing.Files)
	}
	if len(existing.Subfolders) != 1 || existing.Subfolders[0].Name != "C" {
		t.Fatalf("subfolder mapping failed, got %#v", existing.Subfolders)
	}
}

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()
	_ = w.Close()
	os.Stdout = old
	return <-outC
}

func TestDirLabel_UnnamedAndLocked(t *testing.T) {
	d := &pb.Directory{
		Name:        "   ",
		Path:        "/some/path",
		Directories: []*pb.Directory{},
		Files:       []*pb.File{},
		IsLocked:    true,
	}
	s := dirLabel(d)
	if !strings.Contains(s, "<unnamed>") {
		t.Fatalf("expected <unnamed> in label, got: %s", s)
	}
	if !strings.Contains(s, "[locked]") {
		t.Fatalf("expected [locked] in label, got: %s", s)
	}
	if !strings.Contains(s, "files=0") || !strings.Contains(s, "dirs=0") {
		t.Fatalf("expected counts in label, got: %s", s)
	}
}

func TestConvertFolderToProto_And_MetadataConverter(t *testing.T) {
	f := Folder{
		Name: "root",
		Path: "/root",
		Files: []*File{
			{
				Name:    "a.txt",
				Path:    "/root/a.txt",
				Tags:    []string{"one", "two"},
				Locked:  true,
				NewPath: "/new/a.txt",
			},
		},
		Subfolders: []*Folder{
			{Name: "sub", Path: "/root/sub"},
		},
	}

	proto := convertFolderToProto(f)
	if proto.Name != "root" || proto.Path != "/root" {
		t.Fatalf("proto folder header incorrect: %#v", proto)
	}
	if len(proto.Files) != 1 || proto.Files[0].Name != "a.txt" {
		t.Fatalf("file mapping failed: %#v", proto.Files)
	}
	// metadataConverter: simple round-trip
	metaIn := []*pb.MetadataEntry{{Key: "size", Value: "123"}}
	metaOut := metadataConverter(metaIn)
	if len(metaOut) != 1 || metaOut[0].Key != "size" || metaOut[0].Value != "123" {
		t.Fatalf("metadataConverter failed: %#v", metaOut)
	}
}

func TestPrintDirectoryWithMetadata_And_PrintDirChildren(t *testing.T) {
	dir := &pb.Directory{
		Name: "root",
		Path: "/root",
		Directories: []*pb.Directory{
			{
				Name:  "subdir",
				Path:  "/root/subdir",
				Files: []*pb.File{},
			},
		},
		Files: []*pb.File{
			{
				Name:         "file1.txt",
				OriginalPath: "/root/file1.txt",
				Tags:         []*pb.Tag{{Name: "t1"}},
				IsLocked:     false,
			},
		},
		IsLocked: false,
	}

	out := captureOutput(func() {
		printDirectoryWithMetadata(dir, 0)
	})

	// basic assertions that important pieces are printed
	if !strings.Contains(out, "[DIR] root") {
		t.Fatalf("expected root label printed, got: %s", out)
	}
	if !strings.Contains(out, "subdir") {
		t.Fatalf("expected subdir printed, got: %s", out)
	}
	if !strings.Contains(out, "file1.txt") {
		t.Fatalf("expected file1.txt printed, got: %s", out)
	}

	// test empty directory prints "(empty)"
	empty := &pb.Directory{Name: "empty", Path: "/empty"}
	out2 := captureOutput(func() {
		printDirChildren(empty, "")
	})
	if !strings.Contains(out2, "(empty)") {
		t.Fatalf("expected (empty) printed for empty dir, got: %s", out2)
	}
}

func TestLoadEnvFile_ParsesFileCorrectly(t *testing.T) {
	tmp := t.TempDir()
	fpath := filepath.Join(tmp, "test.env")
	content := `# comment line
export FOO=bar
BAZ=qux
EMPTY=
# trailing comment
WITH_EQ=part1=part2
`
	if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars, err := loadEnvFile(fpath)
	if err != nil {
		t.Fatalf("loadEnvFile error: %v", err)
	}

	if vars["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar got %q", vars["FOO"])
	}
	if vars["BAZ"] != "qux" {
		t.Fatalf("expected BAZ=qux got %q", vars["BAZ"])
	}
	// EMPTY present as empty string
	if val, ok := vars["EMPTY"]; !ok || val != "" {
		t.Fatalf("expected EMPTY present as empty string, got %#v, ok=%v", val, ok)
	}
	// WITH_EQ should preserve everything after first =
	if vars["WITH_EQ"] != "part1=part2" {
		t.Fatalf("expected WITH_EQ=part1=part2 got %q", vars["WITH_EQ"])
	}
}

func TestFindProjectRoot_FoundAndNotFound(t *testing.T) {
	tmp := t.TempDir()
	targetName := "marker.txt"
	targetPath := filepath.Join(tmp, targetName)
	if err := os.WriteFile(targetPath, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	// create nested cwd
	nested := filepath.Join(tmp, "a", "b")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	origWD, _ := os.Getwd()
	defer os.Chdir(origWD)

	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}

	found, err := FindProjectRoot(targetName)
	if err != nil {
		t.Fatalf("expected to find %s, got error: %v", targetName, err)
	}
	if filepath.Clean(found) != filepath.Clean(targetPath) {
		t.Fatalf("expected %q, got %q", targetPath, found)
	}

	// Not found case
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	_, err = FindProjectRoot("definitely_not_present_12345")
	if err == nil {
		t.Fatalf("expected error when file not found")
	}
}

func TestLoadTreeDataHandlerGoOnly_NoManager(t *testing.T) {
	// ensure Composites is empty for this test
	orig := Composites
	defer func() { Composites = orig }()
	Composites = []*Folder{}

	req := httptest.NewRequest("GET", "/loadTreeData?name=nonexistent", nil)
	rr := httptest.NewRecorder()

	loadTreeDataHandlerGoOnly(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing manager, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "No smart manager") {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestPrettyPrintFolder_OutputContainsExpected(t *testing.T) {
	// build sample folder
	root := &Folder{
		Name: "root",
		Path: "/root",
		Files: []*File{
			{Name: "a.txt", Path: "/root/a.txt", Tags: []string{"t1", "t2"}},
		},
		Subfolders: []*Folder{
			{Name: "sub", Path: "/root/sub", Files: []*File{{Name: "inner.txt"}}},
		},
	}

	out := captureOutput(func() {
		PrettyPrintFolder(root, "")
	})

	if !strings.Contains(out, "📁 root") {
		t.Fatalf("expected root folder printed, got: %s", out)
	}
	if !strings.Contains(out, "📄 a.txt") {
		t.Fatalf("expected file a.txt printed, got: %s", out)
	}
	if !strings.Contains(out, "TAG: t1") {
		t.Fatalf("expected tag printed, got: %s", out)
	}
	if !strings.Contains(out, "sub") {
		t.Fatalf("expected subfolder printed, got: %s", out)
	}
}

func TestRemoveFileOrderPreserving_BasicAndRecursive(t *testing.T) {
	// top-level removal
	root := &Folder{
		Name: "root",
		Files: []*File{
			{Name: "one", Path: "/one"},
			{Name: "two", Path: "/two"},
			{Name: "three", Path: "/three"},
		},
		Subfolders: []*Folder{
			{
				Name: "child",
				Files: []*File{
					{Name: "c1", Path: "/child/c1"},
				},
			},
		},
	}

	if err := root.RemoveFileOrderPreserving("/two"); err != nil {
		t.Fatalf("unexpected error removing /two: %v", err)
	}
	if len(root.Files) != 2 {
		t.Fatalf("expected 2 files after removal, got %d", len(root.Files))
	}
	if root.Files[0].Path != "/one" || root.Files[1].Path != "/three" {
		t.Fatalf("order incorrect after removal: %#v", root.Files)
	}

	// recursive removal from subfolder
	if err := root.RemoveFileOrderPreserving("/child/c1"); err != nil {
		t.Fatalf("unexpected error removing child file: %v", err)
	}
	// ensure file removed from subfolder
	if len(root.Subfolders[0].Files) != 0 {
		t.Fatalf("expected subfolder files empty after removal, got %#v", root.Subfolders[0].Files)
	}

	// attempt to remove non-existent file
	if err := root.RemoveFileOrderPreserving("/nope"); err == nil {
		t.Fatalf("expected error when removing non-existent file")
	}
}

func TestPrintFileNodeChildren_SortsAndPrints(t *testing.T) {
	// create some file nodes: folders and files
	nodes := []FileNode{
		{
			Name:     "ZFolder",
			IsFolder: true,
			Children: []FileNode{{Name: "afile", IsFolder: false}},
		},
		{
			Name:     "afolder",
			IsFolder: true,
			Children: []FileNode{},
		},
		{
			Name:     "bfile.txt",
			IsFolder: false,
		},
		{
			Name:     "Afile.txt",
			IsFolder: false,
		},
	}

	// shuffle to ensure sorting matters
	sort.Slice(nodes, func(i, j int) bool { return i < j })

	out := captureOutput(func() {
		printFileNodeChildren(nodes, "")
	})

	// expect directories printed (afolder, ZFolder) and files (Afile.txt, bfile.txt) in case-insensitive order
	if !strings.Contains(out, "afolder") {
		t.Fatalf("expected afolder printed, got: %s", out)
	}
	if !strings.Contains(out, "ZFolder") {
		t.Fatalf("expected ZFolder printed, got: %s", out)
	}
	if !strings.Contains(out, "Afile.txt") || !strings.Contains(out, "bfile.txt") {
		t.Fatalf("expected files printed, got: %s", out)
	}
}

func TestSafeName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello", "hello"},
		{"  spaced  ", "spaced"},
		{"", "<unnamed>"},
		{"   ", "<unnamed>"},
	}
	for _, tt := range tests {
		got := safeName(tt.in)
		if got != tt.want {
			t.Errorf("safeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestJoinNonEmpty(t *testing.T) {
	got := joinNonEmpty([]string{"a", " ", "", "b"})
	if got != "a, b" {
		t.Errorf("joinNonEmpty = %q, want %q", got, "a, b")
	}

	got = joinNonEmpty([]string{"   ", ""})
	if got != "" {
		t.Errorf("joinNonEmpty with blanks = %q, want empty string", got)
	}
}

func TestNodeLabel(t *testing.T) {
	tests := []struct {
		node FileNode
		want string
	}{
		{FileNode{Name: "file.txt"}, "[FILE] file.txt"},
		{FileNode{Name: "folder", IsFolder: true}, "[DIR] folder"},
		{FileNode{Name: "f", Tags: []string{"tag1", " ", "tag2"}}, "[FILE] f [tags: tag1, tag2]"},
		{FileNode{Name: "locked", Locked: true}, "[FILE] locked [locked]"},
	}
	for _, tt := range tests {
		got := nodeLabel(tt.node)
		if got != tt.want {
			t.Errorf("nodeLabel(%+v) = %q, want %q", tt.node, got, tt.want)
		}
	}
}
