package filesystem

//run test: go test -v ./filesystem
import (
	"testing"
)

func TestBasicFileTagMetadata(t *testing.T) {
	file := &File{
		Name: "doc.txt",
		Path: "/docs/doc.txt",
	}

	file.Metadata = append(file.Metadata, &MetadataEntry{Key: "Author", Value: "Alice"})
	file.Tags = append(file.Tags, &Tag{ID: "t1", Name: "Important"})

	if len(file.Metadata) != 1 || file.Metadata[0].Key != "Author" {
		t.Errorf("Metadata not correctly added to file")
	}
	if len(file.Tags) != 1 || file.Tags[0].Name != "Important" {
		t.Errorf("Tag not correctly added to file")
	}
}

func TestAddAndGetFile(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	file := &File{Name: "report.pdf", Path: "/report.pdf"}
	root.AddFile(file)

	retrieved := root.GetFile("/report.pdf")
	if retrieved == nil || retrieved.Name != "report.pdf" {
		t.Errorf("Expected to retrieve 'report.pdf', got nil or wrong file")
	}
}

func TestRecursiveFileRetrieval(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	sub := &Folder{Name: "sub", Path: "/sub"}
	file := &File{Name: "deep.txt", Path: "/sub/deep.txt"}

	sub.AddFile(file)
	root.AddSubfolder(sub)

	found := root.GetFile("/sub/deep.txt")
	if found == nil || found.Name != "deep.txt" {
		t.Errorf("Expected to find 'deep.txt', got nil or wrong file")
	}
}

func TestAddAndGetSubfolder(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	child := &Folder{Name: "child", Path: "/child"}
	root.AddSubfolder(child)

	got := root.GetSubfolder("/child")
	if got == nil || got.Name != "child" {
		t.Errorf("Expected to get subfolder '/child'")
	}
}

func TestNestedSubfolderRetrieval(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	sub1 := &Folder{Name: "sub1", Path: "/sub1"}
	sub2 := &Folder{Name: "sub2", Path: "/sub1/sub2"}

	sub1.AddSubfolder(sub2)
	root.AddSubfolder(sub1)

	found := root.GetSubfolder("/sub1/sub2")
	if found == nil || found.Name != "sub2" {
		t.Errorf("Expected to retrieve '/sub1/sub2', got nil or wrong folder")
	}
}

func TestRemoveFile(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	file := &File{Name: "delete.txt", Path: "/delete.txt"}
	root.AddFile(file)

	if !root.RemoveFile("/delete.txt") {
		t.Error("Expected RemoveFile to succeed")
	}
	if root.GetFile("/delete.txt") != nil {
		t.Error("File was not removed properly")
	}
}

func TestRemoveSubfolder(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	sub := &Folder{Name: "temp", Path: "/temp"}
	root.AddSubfolder(sub)

	if !root.RemoveSubfolder("/temp") {
		t.Error("Expected RemoveSubfolder to succeed")
	}
	if root.GetSubfolder("/temp") != nil {
		t.Error("Subfolder was not removed properly")
	}
}

func TestAddTagToFile(t *testing.T) {
	root := &Folder{Name: "root", Path: "/"}
	file := &File{Name: "task.txt", Path: "/task.txt"}
	root.AddFile(file)

	ok := root.AddTagToFile("/task.txt", "id42", "todo")
	if !ok {
		t.Error("AddTagToFile returned false")
	}
	if len(file.Tags) != 1 || file.Tags[0].Name != "todo" {
		t.Error("Tag not correctly added to file")
	}
}

func TestAddTagToSelf(t *testing.T) {
	f := &Folder{Name: "folder", Path: "/folder"}
	f.AddTagToSelf("t2", "work")

	if len(f.Tags) != 1 || f.Tags[0].Name != "work" {
		t.Error("Tag not correctly added to folder")
	}
}

func TestDisplayOutput(t *testing.T) {
	root := &Folder{Name: "root", Path: "/", Tags: []*Tag{{ID: "root", Name: "main"}}}
	file := &File{
		Name:     "file.txt",
		Path:     "/file.txt",
		Tags:     []*Tag{{ID: "t1", Name: "tagged"}},
		Metadata: []*MetadataEntry{{Key: "size", Value: "123KB"}},
	}
	root.AddFile(file)
	root.Display(0)
}
