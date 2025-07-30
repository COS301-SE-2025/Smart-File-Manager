// filesystem/filesystem_test.go
package filesystem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// helper to reset globals between tests
func resetState(t *testing.T, tempDir string) {
	Composites = nil
	SetManagersFilePath(filepath.Join(tempDir, "main.json"))
}

// Test loadManagerRecords when file does not exist
func TestLoadManagerRecords_NoFile(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)

	recs, err := loadManagerRecords()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if recs != nil {
		t.Fatalf("expected nil slice when file missing, got %v", recs)
	}
}

// Test loadManagerRecords with invalid JSON
func TestLoadManagerRecords_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)

	// write invalid JSON
	if err := os.WriteFile(managersFilePath, []byte(`{bad json`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := loadManagerRecords()
	if err == nil {
		t.Fatal("expected JSON unmarshal error, got nil")
	}
}

// Test loadManagerRecords with valid JSON
func TestLoadManagerRecords_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)

	want := []ManagerRecord{
		{Name: "one", Path: "/tmp/one"},
		{Name: "two", Path: "/tmp/two"},
	}
	data, _ := json.Marshal(want)
	if err := os.MkdirAll(filepath.Dir(managersFilePath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(managersFilePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := loadManagerRecords()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("loadManagerRecords = %v; want %v", got, want)
	}
}

// Test saveManagerRecords writes a well-formatted JSON file
func TestSaveManagerRecords(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)

	os.RemoveAll(filepath.Dir(managersFilePath))
	recs := []ManagerRecord{{Name: "A", Path: "/p"}}
	if err := saveManagerRecords(recs); err != nil {
		t.Fatalf("saveManagerRecords failed: %v", err)
	}

	data, err := os.ReadFile(managersFilePath)
	if err != nil {
		t.Fatal(err)
	}
	var loaded []ManagerRecord
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("written file not valid JSON: %v", err)
	}
	if len(loaded) != 1 || loaded[0] != recs[0] {
		t.Fatalf("saved records %v; want %v", loaded, recs)
	}
}

// Test AddManager appends to Composites and persists file
func TestAddAndRemoveManager(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)

	// create a real dir so ConvertToObject won't error
	tmp := t.TempDir()
	// Add first
	if err := AddManager("X", tmp); err != nil {
		t.Fatalf("AddManager failed: %v", err)
	}
	if len(Composites) != 1 {
		t.Fatalf("expected 1 composite, got %d", len(Composites))
	}

	// Add second
	tmp2 := t.TempDir()
	if err := AddManager("Y", tmp2); err != nil {
		t.Fatalf("AddManager failed: %v", err)
	}
	if len(Composites) != 2 {
		t.Fatalf("expected 2 Composites, got %d", len(Composites))
	}

	// Remove one
	if err := RemoveManager(tmp); err != nil {
		t.Fatalf("RemoveManager failed: %v", err)
	}
	if len(Composites) != 1 || Composites[0].Path != tmp2 {
		t.Fatalf("after RemoveManager, expected only %v, got %v", tmp2, Composites)
	}

	// Check file reflects only tmp2
	recs, err := loadManagerRecords()
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 1 || recs[0].Path != tmp2 {
		t.Fatalf("file has %v; want only %v", recs, tmp2)
	}
}

// Test startUpHandler happy path
func TestStartUpHandler_Success(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)
	Composites = nil

	// Add two managers
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()
	if err := AddManager("M1", tmp1); err != nil {
		t.Fatal(err)
	}
	if err := AddManager("M2", tmp2); err != nil {
		t.Fatal(err)
	}

	// reset Composites so startUpHandler rebuilds from disk
	Composites = nil

	req := httptest.NewRequest(http.MethodGet, "/startUp", nil)
	rr := httptest.NewRecorder()
	startUpHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned status %d; want %d", status, http.StatusOK)
	}

	var resp startUpResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	expectedMsg := "Request successful!, Composites: " + strconv.Itoa(2)
	if resp.ResponseMessage != expectedMsg {
		t.Errorf("got message %q; want %q", resp.ResponseMessage, expectedMsg)
	}
	if len(resp.ManagerNames) != 2 ||
		resp.ManagerNames[0] != "M1" ||
		resp.ManagerNames[1] != "M2" {
		t.Errorf("got names %v; want [M1 M2]", resp.ManagerNames)
	}
}

// Test startUpHandler when ReadFile fails
func TestStartUpHandler_LoadError(t *testing.T) {
	// point the managersFilePath at a directory
	dir := t.TempDir()
	SetManagersFilePath(dir)
	Composites = nil

	req := httptest.NewRequest(http.MethodGet, "/startUp", nil)
	rr := httptest.NewRecorder()
	startUpHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusBadRequest)
	}
	body := rr.Body.String()
	if !bytes.Contains([]byte(body), []byte("Internal error:")) {
		t.Errorf("body = %q; want to contain Internal error", body)
	}
}

// replace your flakyWriter definition with this:

type flakyWriter struct {
	header     http.Header
	status     int
	body       bytes.Buffer
	writeCount int
}

// Header must be a method, not a field.
func (w *flakyWriter) Header() http.Header {
	return w.header
}

func (w *flakyWriter) WriteHeader(code int) {
	w.status = code
}

func (w *flakyWriter) Write(p []byte) (int, error) {
	w.writeCount++
	if w.writeCount == 1 {
		// simulate a failure on the first write
		return 0, fmt.Errorf("flaky write error")
	}
	return w.body.Write(p)
}

// then in your test init it like this:
func TestStartUpHandler_EncodeError(t *testing.T) {
	dir := t.TempDir()
	resetState(t, dir)
	Composites = nil

	// write an empty slice so loadManagerRecords passes
	if err := saveManagerRecords([]ManagerRecord{}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/startUp", nil)
	w := &flakyWriter{header: make(http.Header)}
	startUpHandler(w, req)

	// after the first JSON write fails, we expect a 500
	if w.status != http.StatusInternalServerError {
		t.Fatalf("status = %d; want %d", w.status, http.StatusInternalServerError)
	}
	if !bytes.Contains(w.body.Bytes(), []byte("Failed to encode response")) {
		t.Errorf("body = %q; want fallback error", w.body.String())
	}
}
