package filesystem

import (
	"testing"
	"time"
)

// helper to create a file
func newFile(name, path string) *File {
	return &File{Name: name, Path: path}
}

// helper to create a folder
func newFolder(name, path string) *Folder {
	return &Folder{Name: name, Path: path, CreationDate: time.Now()}
}

func TestAddAndGetFile(t *testing.T) {
	r := newFolder("root", "/root")
	f1 := newFile("file1.txt", "/root/file1.txt")
	r.AddFile(f1)
	got := r.GetFile("/root/file1.txt")
	if got == nil {
		t.Fatalf("expected to find file, got nil")
	}
	if got.Name != "file1.txt" {
		t.Errorf("expected Name 'file1.txt', got '%s'", got.Name)
	}
}

func TestRemoveFile(t *testing.T) {
	r := newFolder("root", "/")
	f := newFile("a.txt", "/a.txt")
	r.AddFile(f)
	// remove existing
	removed := r.RemoveFile("/a.txt")
	if !removed {
		t.Errorf("RemoveFile returned false for existing file")
	}
	// ensure it's gone
	if r.GetFile("/a.txt") != nil {
		t.Errorf("file still found after removal")
	}
	// remove non-existing
	if r.RemoveFile("/nonexistent.txt") {
		t.Errorf("RemoveFile returned true for non-existing file")
	}
}

func TestAddAndGetSubfolder(t *testing.T) {
	r := newFolder("root", "/")
	sub := newFolder("sub", "/sub")
	r.AddSubfolder(sub)
	// direct get
	if got := r.GetSubfolder("/sub"); got == nil || got.Name != "sub" {
		t.Fatalf("expected to get sub folder, got %v", got)
	}
	// nested get
	nested := newFolder("nested", "/sub/nested")
	sub.AddSubfolder(nested)
	if got := r.GetSubfolder("/sub/nested"); got == nil || got.Name != "nested" {
		t.Errorf("expected nested folder, got %v", got)
	}
}

func TestRemoveSubfolder(t *testing.T) {
	r := newFolder("root", "/")
	sub := newFolder("x", "/x")
	r.AddSubfolder(sub)
	if !r.RemoveSubfolder("/x") {
		t.Errorf("RemoveSubfolder failed on existing folder")
	}
	if r.GetSubfolder("/x") != nil {
		t.Errorf("subfolder still found after removal")
	}
	// non existing
	if r.RemoveSubfolder("/y") {
		t.Errorf("RemoveSubfolder returned true for non-existing folder")
	}
}

func TestTagging(t *testing.T) {
	r := newFolder("root", "/")
	f := newFile("file.txt", "/file.txt")
	r.AddFile(f)

	// tag file
	if !r.AddTagToFile("/file.txt", "important") {
		t.Fatalf("AddTagToFile returned false")
	}
	if len(f.Tags) != 1 || f.Tags[0] != "important" {
		t.Errorf("expected tag 'important' on file, got %v", f.Tags)
	}

	// tag folder
	r.AddTagToSelf("t1", "projects")
	if len(r.Tags) != 1 || r.Tags[0] != "projects" {
		t.Errorf("expected tag 'projects' on folder, got %v", r.Tags)
	}

	// tagging non-existent file
	if r.AddTagToFile("/no.txt", "none") {
		t.Errorf("AddTagToFile returned true for non-existent file")
	}
}

func TestRemoveTag(t *testing.T) {
	r := newFolder("root", "/")
	f := newFile("f.txt", "/f.txt")
	r.AddFile(f)
	r.AddTagToFile("/f.txt", "temp")
	f.Tags = append(f.Tags, "extra")

	success := f.RemoveTag("temp")
	if !success || len(f.Tags) != 1 || f.Tags[0] != "extra" {
		t.Errorf("expected only 'extra' tag left, got %v", f.Tags)
	}

	success = r.RemoveTag("projects") // not present
	if success {
		t.Errorf("expected false removing nonexistent folder tag")
	}
	r.AddTagToSelf("", "cleanup")
	success = r.RemoveTag("cleanup")
	if !success || len(r.Tags) != 0 {
		t.Errorf("expected no folder tags remaining, got %v", r.Tags)
	}
}
