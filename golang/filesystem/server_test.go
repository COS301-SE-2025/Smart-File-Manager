package filesystem

import (
	"net/http/httptest"
	"testing"
)

func TestAPI_AddAndRemoveDirectory(t *testing.T) {
	req := httptest.NewRequest("GET", "/addDirectory?name=testdir&path=./testdata", nil)
	w := httptest.NewRecorder()
	addCompositeHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("expected true, got %s", w.Body.String())
	}

	req = httptest.NewRequest("GET", "/removeDirectory?path=./testdata", nil)
	w = httptest.NewRecorder()
	removeCompositeHandler(w, req)
	if w.Body.String() != "true" {
		t.Fatalf("expected true, got %s", w.Body.String())
	}
}
