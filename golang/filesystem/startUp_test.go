package filesystem

//generated tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

// setupTestStorage redirects the on-disk file to a temp location
func setupTestStorage(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "main.json")
	SetManagersFilePath(file)
	return func() { os.RemoveAll(dir) }
}

func TestLoadAndSaveManagerRecords(t *testing.T) {
	cleanup := setupTestStorage(t)
	defer cleanup()

	// prepare some records
	want := []ManagerRecord{
		{Name: "A", Path: "/path/A"},
		{Name: "B", Path: "/path/B"},
	}

	// save them
	if err := saveManagerRecords(want); err != nil {
		t.Fatalf("saveManagerRecords error: %v", err)
	}

	// load them back
	got, err := loadManagerRecords()
	if err != nil {
		t.Fatalf("loadManagerRecords error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loadManagerRecords = %v; want %v", got, want)
	}
}

func TestAddAndRemoveManager(t *testing.T) {
	cleanup := setupTestStorage(t)
	defer cleanup()

	// clear any in-memory composites
	composites = nil

	// create two real temp dirs for ConvertToObject
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	// add managers
	if err := AddManager("M1", dir1); err != nil {
		t.Fatalf("AddManager M1 error: %v", err)
	}
	if err := AddManager("M2", dir2); err != nil {
		t.Fatalf("AddManager M2 error: %v", err)
	}

	// in-memory slice should have two entries
	if len(composites) != 2 {
		t.Fatalf("composites length = %d; want 2", len(composites))
	}

	// on-disk should also have 2
	recs, err := loadManagerRecords()
	if err != nil {
		t.Fatalf("loadManagerRecords after adds: %v", err)
	}
	if len(recs) != 2 {
		t.Errorf("records length = %d; want 2", len(recs))
	}

	// remove first
	if err := RemoveManager(dir1); err != nil {
		t.Fatalf("RemoveManager error: %v", err)
	}

	// now in-memory should have just one, with Path=dir2
	if len(composites) != 1 || composites[0].Path != dir2 {
		t.Errorf("composites after remove = %v; want only [%q]", composites, dir2)
	}

	// on-disk should also reflect one
	recs, err = loadManagerRecords()
	if err != nil {
		t.Fatalf("loadManagerRecords after remove: %v", err)
	}
	if len(recs) != 1 || recs[0].Path != dir2 {
		t.Errorf("records after remove = %v; want only [%q]", recs, dir2)
	}
}

func TestStartUpHandler(t *testing.T) {
	cleanup := setupTestStorage(t)
	defer cleanup()

	// clear in-memory
	composites = nil

	// 1) No file on disk → composites stays empty
	req := httptest.NewRequest("GET", "/startup", nil)
	w := httptest.NewRecorder()
	startUpHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d; want 200", res.StatusCode)
	}
	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)
	if !bytes.Contains(body.Bytes(), []byte(`composites: 0`)) {
		t.Errorf("body = %q; want composites: 0", body.String())
	}
	if len(composites) != 0 {
		t.Errorf("composites = %v; want empty", composites)
	}

	// 2) Write two records to disk
	temp1 := t.TempDir()
	temp2 := t.TempDir()
	wantRecs := []ManagerRecord{
		{Name: "X", Path: temp1},
		{Name: "Y", Path: temp2},
	}
	if err := saveManagerRecords(wantRecs); err != nil {
		t.Fatalf("saveManagerRecords: %v", err)
	}

	// clear in-memory and call handler again
	composites = nil
	req = httptest.NewRequest("GET", "/startup", nil)
	w = httptest.NewRecorder()
	startUpHandler(w, req)
	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status2 = %d; want 200", res.StatusCode)
	}
	body.Reset()
	body.ReadFrom(res.Body)
	if !bytes.Contains(body.Bytes(), []byte(`composites: 2`)) {
		t.Errorf("body2 = %q; want composites: 2", body.String())
	}
	// in-memory slice should now have two entries
	if len(composites) != 2 {
		t.Errorf("composites after startup = %v; want 2 entries", composites)
	}

	// Optional: check that the printed message is valid JSON and contains message
	var msg struct{ Message string }
	if err := json.Unmarshal(body.Bytes(), &msg); err != nil {
		t.Errorf("response not valid JSON: %v", err)
	}
	if !bytes.Contains([]byte(msg.Message), []byte("composites: "+strconv.Itoa(len(wantRecs)))) {
		t.Errorf("message = %q; want composites count", msg.Message)
	}
}
