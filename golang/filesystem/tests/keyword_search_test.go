package test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
	"github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

// Helper: create a temp text file with given content.
func writeTempTextFile(
	t *testing.T,
	dir string,
	name string,
	content string,
) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return p
}

func TestLevenshteinDistForKeywords(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		if got := filesystem.LevenshteinDistForKeywords("", ""); got != 0 {
			t.Fatalf("want 0, got %d", got)
		}
	})

	t.Run("exact match case-insensitive", func(t *testing.T) {
		if got := filesystem.LevenshteinDistForKeywords("Hello", "hello"); got != 0 {
			t.Fatalf("want 0, got %d", got)
		}
	})

	t.Run("substring boost", func(t *testing.T) {
		got := filesystem.LevenshteinDistForKeywords("learn", "deep learning")
		if got > 1 {
			t.Fatalf("want <= 1 due to boost, got %d", got)
		}
	})

	t.Run("far apart", func(t *testing.T) {
		got := filesystem.LevenshteinDistForKeywords("cat", "hippopotamus")
		if got <= 3 {
			t.Fatalf("want > 3, got %d", got)
		}
	})
}

func TestExtractKeywordsFromText_Basic(t *testing.T) {
	// Use a common sample phrase; allow more results to reduce tie flakiness.
	text := "The quick brown fox jumps over the lazy dog. Quick fox quick."
	kws := filesystem.ExtractKeywordsFromText(text, 10)

	if len(kws) == 0 {
		t.Fatalf("want some keywords, got 0")
	}

	// Ensure stopwords are removed
	for _, k := range kws {
		if k.Keyword == "the" || k.Keyword == "is" || k.Keyword == "and" {
			t.Fatalf("unexpected stopword in keywords: %q", k.Keyword)
		}
	}

	// Accept any of these high-signal terms; RAKE ranking may vary on ties.
	wantAny := map[string]struct{}{
		"quick": {}, "fox": {}, "brown": {}, "jumps": {}, "over": {},
	}
	foundAny := false
	for _, k := range kws {
		if _, ok := wantAny[k.Keyword]; ok {
			foundAny = true
			break
		}
	}
	if !foundAny {
		t.Fatalf("expected at least one of quick/fox/brown/jumps/over in keywords, got: %+v", kws)
	}
}

func TestExtractKeywordsRAKE_TextAndBinaryAndSize(t *testing.T) {
	dir := t.TempDir()

	// Text file
	textPath := writeTempTextFile(
		t,
		dir,
		"doc.txt",
		"Go is fun. Go concurrency is powerful. Concurrency patterns!",
	)
	kws, err := filesystem.ExtractKeywordsRAKE(textPath, 10, 1<<20)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(kws) == 0 {
		t.Fatalf("expected some keywords from text file")
	}

	// Binary file (PNG header)
	binPath := filepath.Join(dir, "img.bin")
	bin := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	if err := os.WriteFile(binPath, bin, 0o644); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	kws, err = filesystem.ExtractKeywordsRAKE(binPath, 10, 1<<20)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(kws) != 0 {
		t.Fatalf("expected no keywords for non-text file, got %d", len(kws))
	}

	// Too large file (size check)
	largePath := writeTempTextFile(
		t,
		dir,
		"large.txt",
		strings.Repeat("A", 2048),
	)
	kws, err = filesystem.ExtractKeywordsRAKE(largePath, 10, 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(kws) != 0 {
		t.Fatalf("expected no keywords due to size limit, got %d", len(kws))
	}
}

func TestKeywordSearchReadyHandler_ExistingFalseAndMissing(t *testing.T) {
	// Arrange composites (hasKeywords defaults to false, which we can't set here).
	comp := &filesystem.Folder{Name: "ReadyComp"}

	// Set package global (restore afterward).
	old := filesystem.Composites
	filesystem.Composites = []*filesystem.Folder{comp}
	t.Cleanup(func() { filesystem.Composites = old })

	// Existing comp should return false (not ready) by default
	{
		q := url.Values{}
		q.Set("compositeName", "ReadyComp")
		req := httptest.NewRequest("GET", "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()

		// Exported handler
		filesystem.IsKeywordSearchReadyHander(rec, req)

		if rec.Code != 200 {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var ready bool
		if err := json.Unmarshal(rec.Body.Bytes(), &ready); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if ready {
			t.Fatalf("expected not ready (false) by default")
		}
	}

	// Missing comp should produce 400 (Bad Request) per current handler
	{
		q := url.Values{}
		q.Set("compositeName", "MissingComp")
		req := httptest.NewRequest("GET", "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()

		filesystem.IsKeywordSearchReadyHander(rec, req)

		if rec.Code != 400 {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	}
}

func TestKeywordSearchHandler_FindsByApproxAndSetsCORS(t *testing.T) {
	// Composite with one file that has keyword "learning"
	comp := &filesystem.Folder{
		Name: "C1",
		Files: []*filesystem.File{
			{
				Name: "file1.txt",
				Path: "/p1",
				Keywords: []*pb.Keyword{
					{Keyword: "learning", Score: 1},
				},
			},
		},
	}

	// Set package global (restore afterward).
	old := filesystem.Composites
	filesystem.Composites = []*filesystem.Folder{comp}
	t.Cleanup(func() { filesystem.Composites = old })

	q := url.Values{}
	q.Set("compositeName", "C1")
	q.Set("searchText", "learn") // approximate to "learning"
	req := httptest.NewRequest("GET", "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()

	filesystem.KeywordSearchHadler(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	// CORS headers present
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("CORS Allow-Origin = %q, want *", got)
	}

	// Decode response
	var resp struct {
		Name     string `json:"name"`
		IsFolder bool   `json:"isFolder"`
		Children []struct {
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"children"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Name != "C1" {
		t.Fatalf("resp.Name = %q, want C1", resp.Name)
	}
	if len(resp.Children) == 0 {
		t.Fatalf("expected at least one matched file")
	}
	if resp.Children[0].Name != "file1.txt" {
		t.Fatalf("top result = %q, want file1.txt", resp.Children[0].Name)
	}
}

func TestKeywordSearchHandler_MissingComposite_BadRequest(t *testing.T) {
	// No composites registered
	old := filesystem.Composites
	filesystem.Composites = nil
	t.Cleanup(func() { filesystem.Composites = old })

	q := url.Values{}
	q.Set("compositeName", "DoesNotExist")
	q.Set("searchText", "anything")
	req := httptest.NewRequest("GET", "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()

	filesystem.KeywordSearchHadler(rec, req)

	// Current handler returns 400 (BadRequest). If you change the handler
	// to use http.StatusNotFound, update this assertion to 404.
	if rec.Code != 400 {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestAppendUniqueKeywords_SmallMerge_UniqueAndSkip(t *testing.T) {
	dst := kws("a", "b", "c", "d", "e")
	src := []*pb.Keyword{
		nil,
		{Keyword: "c"}, // dup
		{Keyword: ""},  // empty -> skip
		{Keyword: "f"},
		{Keyword: "a"}, // dup
		{Keyword: "g"},
	}

	got := filesystem.AppendUniqueKeywords(dst, src)
	want := []string{"a", "b", "c", "d", "e", "f", "g"}

	mustEqual(t, namesOf(got), want)
}

func TestAppendUniqueKeywords_LargeMerge_UniqueAndOrder(t *testing.T) {
	dst := kws("k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8", "k9", "k10")
	src := kws("k5", "k11", "k12", "k1", "k13", "k14") // dups: k5, k1

	got := filesystem.AppendUniqueKeywords(dst, src)
	want := []string{
		"k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8", "k9", "k10",
		"k11", "k12", "k13", "k14",
	}

	mustEqual(t, namesOf(got), want)
}

func TestAppendUniqueKeywords_CapAt30_WhenAppending(t *testing.T) {
	// dst has 28, src has 4 (2 unique) -> result should be 30 with s1, s2
	dst := seq("d", 28) // d1..d28
	src := kws("d28", "s1", "s2", "s3", "s4")

	got := filesystem.AppendUniqueKeywords(dst, src)

	want := append(namesOf(seq("d", 28)), "s1", "s2")
	if len(got) != 30 {
		t.Fatalf("cap not applied: got len=%d, want len=30", len(got))
	}
	mustEqual(t, namesOf(got), want)
}

func TestAppendUniqueKeywords_DstExactlyAtCap_Unchanged(t *testing.T) {
	dst := seq("x", 30)
	base := namesOf(dst)

	got := filesystem.AppendUniqueKeywords(dst, kws("new1", "new2"))
	mustEqual(t, namesOf(got), base)
}

func TestAppendUniqueKeywords_DstOverCap_IsTrimmedTo30(t *testing.T) {
	dst := seq("x", 35)
	got := filesystem.AppendUniqueKeywords(dst, kws("new1", "new2"))

	if len(got) != 30 {
		t.Fatalf("expected trim to 30, got len=%d", len(got))
	}
	want := namesOf(seq("x", 30))
	mustEqual(t, namesOf(got), want)
}

func TestAppendUniqueKeywords_CaseSensitive(t *testing.T) {
	dst := kws("Dog")
	src := kws("dog") // case-sensitive -> should append

	got := filesystem.AppendUniqueKeywords(dst, src)
	want := []string{"Dog", "dog"}

	mustEqual(t, namesOf(got), want)
}

func TestAppendUniqueKeywords_EmptySrc_NoChange(t *testing.T) {
	dst := kws("a", "b")
	base := namesOf(dst)

	got := filesystem.AppendUniqueKeywords(dst, nil)
	mustEqual(t, namesOf(got), base)

	got = filesystem.AppendUniqueKeywords(dst, []*pb.Keyword{})
	mustEqual(t, namesOf(got), base)
}

// Helpers

func kws(ss ...string) []*pb.Keyword {
	out := make([]*pb.Keyword, 0, len(ss))
	for _, s := range ss {
		out = append(out, &pb.Keyword{Keyword: s})
	}
	return out
}

func seq(prefix string, n int) []*pb.Keyword {
	out := make([]*pb.Keyword, 0, n)
	for i := 1; i <= n; i++ {
		out = append(out, &pb.Keyword{Keyword: fmt.Sprintf("%s%d", prefix, i)})
	}
	return out
}

func namesOf(kws []*pb.Keyword) []string {
	out := make([]string, 0, len(kws))
	for _, k := range kws {
		if k == nil {
			out = append(out, "<nil>")
			continue
		}
		out = append(out, k.Keyword)
	}
	return out
}

func mustEqual(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mismatch:\n got:  %#v\n want: %#v", got, want)
	}
}
