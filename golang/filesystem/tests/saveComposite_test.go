package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
	"github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

func TestSavePopulateDeleteCompositeDetails_RoundTrip(t *testing.T) {
	// Work in a temp dir so storage/ is isolated.
	tmp := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	// Build a composite with files and a nested subfolder.
	f1 := &filesystem.File{
		Name:   "f1.txt",
		Path:   filepath.Join(tmp, "f1.txt"),
		Tags:   []string{"t1", "t2"},
		Locked: true,
		Keywords: []*pb.Keyword{
			{Keyword: "alpha", Score: 1.0},
			{Keyword: "beta", Score: 0.8},
		},
	}
	f2 := &filesystem.File{
		Name:   "f2.txt",
		Path:   filepath.Join(tmp, "f2.txt"),
		Tags:   []string{"x"},
		Locked: false,
		Keywords: []*pb.Keyword{
			{Keyword: "gamma", Score: 0.7},
		},
	}
	subFile := &filesystem.File{
		Name:   "sf1.md",
		Path:   filepath.Join(tmp, "sub", "sf1.md"),
		Tags:   []string{"subtag"},
		Locked: true,
		Keywords: []*pb.Keyword{
			{Keyword: "delta", Score: 0.9},
		},
	}

	comp := &filesystem.Folder{
		Name: "compA",
		Path: tmp,
		Files: []*filesystem.File{
			f1, f2,
		},
		Subfolders: []*filesystem.Folder{
			{
				Name:  "sub",
				Path:  filepath.Join(tmp, "sub"),
				Files: []*filesystem.File{subFile},
			},
		},
	}

	// Save JSON to storage/compA.json (via test-only wrapper).
	filesystem.SaveCompositeDetailsForTest(comp)

	storageDir := "storage"
	jsonPath := filepath.Join(storageDir, comp.Name+".json")

	// Ensure only the final JSON exists (no tmp-*.json left behind).
	ents, err := os.ReadDir(storageDir)
	if err != nil {
		t.Fatalf("ReadDir(storage): %v", err)
	}
	if len(ents) != 1 || ents[0].Name() != comp.Name+".json" {
		t.Fatalf("expected only %s in storage, got %v", comp.Name+".json", ents)
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("expected saved JSON at %s: %v", jsonPath, err)
	}

	// Inspect saved structure.
	var tree filesystem.DirectoryTreeJson
	if err := json.Unmarshal(data, &tree); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tree.Name != comp.Name || !tree.IsFolder || tree.RootPath != comp.Path {
		t.Fatalf("root mismatch: %+v", tree)
	}

	// Collect top-level file nodes and the subfolder node.
	filesByPath := map[string]filesystem.FileNode{}
	var subNode *filesystem.FileNode
	for i := range tree.Children {
		n := tree.Children[i]
		if !n.IsFolder {
			filesByPath[n.Path] = n
		} else if n.Name == "sub" {
			subNode = &n
		}
	}
	if _, ok := filesByPath[f1.Path]; !ok {
		t.Fatalf("missing f1 in JSON")
	}
	if _, ok := filesByPath[f2.Path]; !ok {
		t.Fatalf("missing f2 in JSON")
	}
	if subNode == nil {
		t.Fatalf("missing subfolder node")
	}
	if len(subNode.Children) != 1 || subNode.Children[0].Path != subFile.Path {
		t.Fatalf("subfolder children mismatch: %+v", subNode.Children)
	}
	// Check keywords preserved for a file
	if len(filesByPath[f1.Path].Keywords) == 0 ||
		filesByPath[f1.Path].Keywords[0].Keyword != "alpha" {
		t.Fatalf("f1 keywords not saved correctly: %+v", filesByPath[f1.Path].Keywords)
	}
	if !filesByPath[f1.Path].Locked || len(filesByPath[f1.Path].Tags) != 2 {
		t.Fatalf("f1 tags/locked not saved: tags=%v locked=%v",
			filesByPath[f1.Path].Tags, filesByPath[f1.Path].Locked,
		)
	}

	// Build a fresh composite (empty metadata) and repopulate from JSON.
	f1b := &filesystem.File{Name: "f1.txt", Path: f1.Path}
	f2b := &filesystem.File{Name: "f2.txt", Path: f2.Path}
	subFileB := &filesystem.File{Name: "sf1.md", Path: subFile.Path}
	comp2 := &filesystem.Folder{
		Name:  "compA", // must match for populateKeywordsFromStoredJsonFile
		Path:  tmp,
		Files: []*filesystem.File{f1b, f2b},
		Subfolders: []*filesystem.Folder{
			{
				Name:  "sub",
				Path:  filepath.Join(tmp, "sub"),
				Files: []*filesystem.File{subFileB},
			},
		},
	}

	filesystem.PopulateKeywordsFromStoredJsonFileForTest(comp2)

	// Top-level files should have keywords, tags, and locked restored.
	if len(f1b.Keywords) == 0 || f1b.Keywords[0].Keyword != "alpha" {
		t.Fatalf("f1 keywords not populated: %+v", f1b.Keywords)
	}
	if got := len(f1b.Tags); got != 2 || !f1b.Locked {
		t.Fatalf("f1 tags/locked not populated: tags=%v locked=%v",
			f1b.Tags, f1b.Locked,
		)
	}
	if len(f2b.Keywords) == 0 || f2b.Keywords[0].Keyword != "gamma" {
		t.Fatalf("f2 keywords not populated: %+v", f2b.Keywords)
	}
	if got := len(f2b.Tags); got != 1 || f2b.Locked {
		t.Fatalf("f2 tags/locked not populated correctly: tags=%v locked=%v",
			f2b.Tags, f2b.Locked,
		)
	}

	// Nested file: helperMerge sets only Keywords (by current implementation).
	if len(subFileB.Keywords) == 0 || subFileB.Keywords[0].Keyword != "delta" {
		t.Fatalf("subfile keywords not populated: %+v", subFileB.Keywords)
	}
	// if len(subFileB.Tags) != 0 || subFileB.Locked {
	// 	t.Fatalf("subfile tags/locked should not be populated by helper; got tags=%v locked=%v",
	// 		subFileB.Tags, subFileB.Locked,
	// 	)
	// }

	// Delete the JSON and ensure idempotency.
	if err := filesystem.DeleteCompositeDetailsFileForTest(comp2.Name); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		t.Fatalf("file should be deleted")
	}
	if err := filesystem.DeleteCompositeDetailsFileForTest(comp2.Name); err != nil {
		t.Fatalf("second delete should be nil, got: %v", err)
	}
}

func TestPopulateKeywordsFromStoredJsonFile_NoFile(t *testing.T) {
	// Work in a temp dir so storage/ is isolated.
	tmp := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	// Composite without a saved JSON present.
	f1 := &filesystem.File{Name: "f1.txt", Path: filepath.Join(tmp, "f1.txt")}
	comp := &filesystem.Folder{
		Name:  "nojson",
		Path:  tmp,
		Files: []*filesystem.File{f1},
	}

	// Should silently do nothing (no panic, no changes).
	filesystem.PopulateKeywordsFromStoredJsonFileForTest(comp)

	if len(f1.Keywords) != 0 || len(f1.Tags) != 0 || f1.Locked {
		t.Fatalf("unexpected changes when no JSON present: keywords=%v tags=%v locked=%v",
			f1.Keywords, f1.Tags, f1.Locked,
		)
	}
}
