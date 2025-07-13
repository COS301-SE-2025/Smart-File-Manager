package filesystem

import (
	"net/http/httptest"
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
