package search_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

func TestLevenshteinDist(t *testing.T) {
	tests := []struct {
		searchText string
		fileName   string
		want       int
		desc       string
	}{
		// Exact match (should return 0)
		{"report", "report", 0, "Exact match"},
		{"Report", "report", 0, "Case-insensitive exact match"},

		// Substring match (should return low score)
		{"jack", "jacks homework.txt", 2, "Substring match boosts score"},

		// Extensions ignored
		{"budget", "budget.pdf", 0, "Extension stripped - exact match"},
		{"budget", "budget_report.pdf", 1, "Extension stripped - partial match"},

		// Completely different strings
		{"dog", "cat", 3, "Different strings"},

		// Prefix match
		{"home", "homework.txt", 1, "Prefix substring boosts"},

		// Suffix match
		{"work", "homework.txt", 1, "Suffix substring boosts"},

		// Shorter query, longer file
		{"a", "abcdefgh.txt", 1, "Very short search in long filename"},

		// Empty string cases
		{"", "anything.txt", 12, "Empty search text"},
		{"query", "", 5, "Empty file name"},
	}

	for _, tt := range tests {
		got := filesystem.LevenshteinDist(tt.searchText, tt.fileName)
		if got != tt.want {
			t.Errorf("FAILED [%s]: LevenshteinDist(%q, %q) = %d; want %d",
				tt.desc, tt.searchText, tt.fileName, got, tt.want)
		}
	}
}

// TestConvertMetadataEntries ensures that keys are matched case-insensitively
// and only the known fields are set.
func TestConvertMetadataEntries(t *testing.T) {
	entries := []*filesystem.MetadataEntry{
		{Key: "Size", Value: "123"},
		{Key: "DateCreated", Value: "2025-07-30 14:00"},
		{Key: "LastModified", Value: "2025-07-30 15:00"},
		{Key: "MimeType", Value: ".txt"},
		// unknown key should be ignored
		{Key: "Unknown", Value: "ignore"},
	}

	md := filesystem.ConvertMetadataEntries(entries)

	want := &filesystem.Metadata{
		Size:         "123",
		DateCreated:  "2025-07-30 14:00",
		LastModified: "2025-07-30 15:00",
		MimeType:     ".txt",
	}
	if !reflect.DeepEqual(md, want) {
		t.Errorf("ConvertMetadataEntries() = %+v; want %+v", md, want)
	}
}

// fakeFolder is a helper to create a simple Folder for testing SearchHandler
func fakeFolder(name string) *filesystem.Folder {
	return &filesystem.Folder{
		Name:       name,
		Path:       "/fake/" + name,
		Files:      []*filesystem.File{},
		Subfolders: []*filesystem.Folder{},
	}
}

// // TestSearchHandlerNotFound tests that an unknown compositeName returns 404
// func TestSearchHandlerNotFound(t *testing.T) {
// 	// clear Composites
// 	filesystem.Composites = []*filesystem.Folder{}

// 	req := httptest.NewRequest("GET", "/search?compositeName=foo&searchText=bar", nil)
// 	rr := httptest.NewRecorder()

// 	filesystem.SearchHandler(rr, req)

// 	if rr.Code != http.StatusNotFound {
// 		t.Errorf("expected status 404; got %d", rr.Code)
// 	}
// }

// TestSearchHandlerEmpty tests that a known composite with no files returns empty children list
func TestSearchHandlerEmpty(t *testing.T) {
	filesystem.Composites = []*filesystem.Folder{fakeFolder("test")}

	req := httptest.NewRequest("GET", "/search?compositeName=test&searchText=anything", nil)
	rr := httptest.NewRecorder()

	filesystem.SearchHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200; got %d", rr.Code)
	}
	if ctype := rr.Header().Get("Content-Type"); !strings.Contains(ctype, "application/json") {
		t.Errorf("expected Content-Type application/json; got %s", ctype)
	}

	// decode body
	var resp struct {
		Name     string        `json:"name"`
		IsFolder bool          `json:"isFolder"`
		Children []interface{} `json:"children"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if resp.Name != "test" {
		t.Errorf("expected Name 'test'; got %q", resp.Name)
	}
	if !resp.IsFolder {
		t.Errorf("expected IsFolder true; got false")
	}
	if len(resp.Children) != 0 {
		t.Errorf("expected empty Children; got %v", resp.Children)
	}
}
