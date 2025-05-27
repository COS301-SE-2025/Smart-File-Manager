package filesystem

//run test: go test -v ./filesystem
import (
	"fmt"
	"testing"
	"time"
)

func createTestStructure() *Folder {
	root := &Folder{managedItem: managedItem{
		ItemPath: "/root", ItemName: "root",
	}}

	sub := &Folder{managedItem: managedItem{
		ItemPath: "/root/sub", ItemName: "sub",
	}}

	file := &File{managedItem: managedItem{
		ItemPath:     "/root/sub/file.txt",
		ItemName:     "file.txt",
		CreationDate: time.Now(),
	}}

	err := sub.AddItem(file)
	if err != nil {
		fmt.Println(err)
	}
	err1 := root.AddItem(sub)
	if err1 != nil {
		fmt.Println(err1)
	}

	return root
}

func TestRemoveFile(t *testing.T) {
	root := createTestStructure()

	removed := root.RemoveItem("/root/sub/file.txt")
	if !removed {
		t.Errorf("Expected file to be removed, but it was not")
	}
}

func TestRemoveNonExistentFile(t *testing.T) {
	root := createTestStructure()

	removed := root.RemoveItem("/root/nope.txt")
	if removed {
		t.Errorf("Expected removal to fail for non-existent file")
	}
}

func TestSubfolderStillExistsAfterFileRemoval(t *testing.T) {
	root := createTestStructure()
	root.RemoveItem("/root/sub/file.txt")

	// Check if subfolder still exists
	found := false
	for _, item := range root.ContainedItems {
		if item.GetPath() == "/root/sub" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected subfolder to still exist after file removal")
	}
}

func TestRecursiveRemoval(t *testing.T) {
	root := createTestStructure()
	sub := root.ContainedItems[0].(*Folder)

	if len(sub.ContainedItems) != 1 {
		t.Fatalf("Expected 1 item in subfolder, got %d", len(sub.ContainedItems))
	}

	root.RemoveItem("/root/sub/file.txt")

	if len(sub.ContainedItems) != 0 {
		t.Errorf("Expected subfolder to be empty after file removal")
	}
}
func TestGetItem(t *testing.T) {
	root := createTestStructure()

	item := root.GetItem("/root/sub/file.txt")
	if item == nil {
		t.Fatalf("Expected to find item at /root/sub/file.txt, got nil")
	}
	if item.GetPath() != "/root/sub/file.txt" {
		t.Errorf("Expected path '/root/sub/file.txt', got '%s'", item.GetPath())
	}
}

func TestAddTagToItem(t *testing.T) {
	root := createTestStructure()

	success := root.AddTagToItem("/root/sub/file.txt", "t1", "Important")
	if !success {
		t.Fatalf("Expected AddTag to succeed, but it failed")
	}

	item := root.GetItem("/root/sub/file.txt")
	if item == nil {
		t.Fatalf("Expected to find item after tagging, got nil")
	}

	file, ok := item.(*File)
	if !ok {
		t.Fatalf("Expected item to be of type *File")
	}

	if len(file.ItemTags) != 1 {
		t.Errorf("Expected 1 tag, found %d", len(file.ItemTags))
	} else if file.ItemTags[0].tagID != "t1" || file.ItemTags[0].tagName != "Important" {
		t.Errorf("Expected tag (t1, Important), got (%s, %s)", file.ItemTags[0].tagID, file.ItemTags[0].tagName)
	}
}
func TestAddTagToNonExistentItem(t *testing.T) {
	root := createTestStructure()

	success := root.AddTagToItem("/root/sub/ghost.txt", "t1", "GhostTag")
	if success {
		t.Errorf("Expected AddTag to fail for non-existent item, but it succeeded")
	}
}

func TestAddTagToFolder(t *testing.T) {
	root := createTestStructure()

	success := root.AddTagToItem("/root/sub", "t2", "ProjectDocs")
	if !success {
		t.Errorf("Expected to successfully add tag to folder")
	}

	// Check if the tag was actually added
	sub := root.GetItem("/root/sub")
	if sub == nil {
		t.Fatalf("Subfolder not found")
	}

	found := false
	if folder, ok := sub.(*Folder); ok {
		for _, tag := range folder.ItemTags {
			if tag.tagID == "t2" && tag.tagName == "ProjectDocs" {
				found = true
			}
		}
	}

	if !found {
		t.Errorf("Expected tag 'ProjectDocs' with ID 't2' to be present in subfolder")
	}
}
