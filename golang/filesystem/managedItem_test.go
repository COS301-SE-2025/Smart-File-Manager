package filesystem

//run test: go test ./filesystem
import (
	"testing"
	"time"
)

func createTestStructure() *Folder {
	root := &Folder{managedItem: managedItem{
		itemPath: "/root", itemName: "root",
	}}

	sub := &Folder{managedItem: managedItem{
		itemPath: "/root/sub", itemName: "sub",
	}}

	file := &File{managedItem: managedItem{
		itemPath:     "/root/sub/file.txt",
		itemName:     "file.txt",
		creationDate: time.Now(),
	}}

	sub.AddItem(file)
	root.AddItem(sub)

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
	for _, item := range root.containedItems {
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
	sub := root.containedItems[0].(*Folder)

	if len(sub.containedItems) != 1 {
		t.Fatalf("Expected 1 item in subfolder, got %d", len(sub.containedItems))
	}

	root.RemoveItem("/root/sub/file.txt")

	if len(sub.containedItems) != 0 {
		t.Errorf("Expected subfolder to be empty after file removal")
	}
}
