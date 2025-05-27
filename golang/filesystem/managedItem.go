package filesystem

import (
	"errors"
	"fmt"
	"strings"
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
	AddItem(item FileSystemItem) error
	GetItem(itemPath string) FileSystemItem
	AddTag(tagID string, tagName string) bool
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

func (f *Folder) AddItem(newItem FileSystemItem) error {
	f.containedItems = append(f.containedItems, newItem)
	return nil
}

func (f *File) AddItem(item FileSystemItem) error {
	return errors.New("cannot add item to a File: operation not supported")
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

func (f *Folder) GetItem(itemPath string) FileSystemItem {
	for _, item := range f.containedItems {
		if item.GetPath() == itemPath {
			return item
		}
		if folder, ok := item.(*Folder); ok {
			if found := folder.GetItem(itemPath); found != nil {
				return found
			}
		}
	}
	return nil
}
func (f *File) GetItem(itemPath string) FileSystemItem {
	//get item needs to be called on folder
	return nil
}

func (f *Folder) AddTagToItem(itemPath string, tagID string, tagName string) bool {
	item := f.GetItem(itemPath)
	if item != nil {
		item.AddTag(tagID, tagName)
		return true
	}
	return false
}

func (f *Folder) AddTagToSelf(tagID string, tagName string) {
	f.itemTags = append(f.itemTags, tag{tagID, tagName})
}

func (f *File) AddTag(tagID string, tagName string) bool {
	f.itemTags = append(f.itemTags, tag{tagID, tagName})
	return true
}

func (f *Folder) AddTag(tagID string, tagName string) bool {
	f.itemTags = append(f.itemTags, tag{tagID, tagName})
	return true
}
func (f *Folder) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFolder: %s\n", prefix, f.itemName)

	if len(f.itemTags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.itemTags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag.tagID, tag.tagName)
		}
	}

	for _, item := range f.containedItems {
		switch v := item.(type) {
		case *Folder:
			v.Display(indent + 1)
		case *File:
			v.Display(indent + 1)
		}
	}
}
func (f *File) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFile: %s\n", prefix, f.itemName)

	if len(f.itemTags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.itemTags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag.tagID, tag.tagName)
		}
	}
}
