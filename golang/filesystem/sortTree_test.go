package filesystem

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestSortTreeHandler_InvalidCase(t *testing.T) {
	req := httptest.NewRequest("GET", "/sortTree?name=TestComp&case=INVALID", nil)
	w := httptest.NewRecorder()

	sortTreeHandler(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid case")
	}
	if !strings.Contains(w.Body.String(), "Invalid case type.") {
		t.Errorf("expected error body, got: %s", w.Body.String())
	}
}

func TestSortTreeHandler_ValidCaseSnake(t *testing.T) {
	mu = sync.Mutex{}
	Composites = []*Folder{} // nothing added

	req := httptest.NewRequest("GET", "/sortTree?name=TestComp&case=SNAKE", nil)
	w := httptest.NewRecorder()

	sortTreeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 BadRequest for missing composite, got %d", resp.StatusCode)
	}
}

func TestSortTreeHandler_NoSuchComposite(t *testing.T) {

	req := httptest.NewRequest("GET", "/sortTree?name=DoesNotExist", nil)
	w := httptest.NewRecorder()

	sortTreeHandler(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for missing composite")
	}
	if !strings.Contains(w.Body.String(), "No smart manager with that name") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}
