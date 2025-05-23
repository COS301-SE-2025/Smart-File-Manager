package filesystem

import (
	"time"
)

// tag structure
type tag struct {
	tagID   string
	tagName string
}

// component interface
type FileSystemItem interface {
	GetPath() string
	RemoveItem(itemPath string) bool
}

// base struct
type managedItem struct {
	itemID       string
	itemName     string
	itemPath     string
	itemTags     []tag
	locked       bool
	fileType     string
	creationDate time.Time
}

func (m *managedItem) GetPath() string {
	return m.itemPath
}

// Leaf
type File struct {
	managedItem
}

func (f *File) RemoveItem(itemPath string) bool {
	// A file has no children; return false
	return false
}

// Composite
type Folder struct {
	managedItem
	containedItems []FileSystemItem
}

func (f *Folder) AddItem(newItem FileSystemItem) {
	f.containedItems = append(f.containedItems, newItem)
}

func (f *Folder) RemoveItem(itemPath string) bool {
	for i, item := range f.containedItems {
		if item.GetPath() == itemPath {
			f.containedItems = append(f.containedItems[:i], f.containedItems[i+1:]...)
			return true
		}
		// if item is a Folder, attempt recursive removal
		if folder, ok := item.(*Folder); ok {
			if folder.RemoveItem(itemPath) {
				return true
			}
		}
	}
	return false
}

func (f *Folder) GetPath() string {
	return f.itemPath
}
