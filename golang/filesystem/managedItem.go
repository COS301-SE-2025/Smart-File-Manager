package filesystem

import (
	"fmt"
	"strings"
	"time"
)

// MetadataEntry holds file metadata key and value
type MetadataEntry struct {
	Key   string
	Value string
}

// File represents a file in the filesystem
type File struct {
	Name     string
	Path     string
	newPath  string
	Metadata []*MetadataEntry
	Tags     []string
	Locked   bool // Lock status for file
}

// Folder represents a directory in the filesystem
type Folder struct {
	Name         string
	Path         string
	newPath      string
	CreationDate time.Time
	Locked       bool // Lock status for folder
	Files        []*File
	Subfolders   []*Folder
	Tags         []string
}

// -------------------- Folder Methods --------------------

// AddFile adds a file to the folder
func (f *Folder) AddFile(file *File) {
	f.Files = append(f.Files, file)
}

// AddSubfolder adds a subfolder to the folder
func (f *Folder) AddSubfolder(folder *Folder) {
	f.Subfolders = append(f.Subfolders, folder)
}

// RemoveFile removes a file by path
func (f *Folder) RemoveFile(filePath string) bool {
	for i, file := range f.Files {
		if file.Path == filePath {
			f.Files = append(f.Files[:i], f.Files[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveSubfolder removes a folder by path
func (f *Folder) RemoveSubfolder(folderPath string) bool {
	for i, folder := range f.Subfolders {
		if folder.Path == folderPath {
			f.Subfolders = append(f.Subfolders[:i], f.Subfolders[i+1:]...)
			return true
		}
	}
	return false
}

// GetFile returns a file by path, searching recursively
func (f *Folder) GetFile(filePath string) *File {
	for _, file := range f.Files {
		if file.Path == filePath {
			return file
		}
	}
	for _, folder := range f.Subfolders {
		if found := folder.GetFile(filePath); found != nil {
			return found
		}
	}
	return nil
}

// GetSubfolder returns a folder by path, searching recursively
func (f *Folder) GetSubfolder(folderPath string) *Folder {
	if f.Path == folderPath {
		return f
	}
	for _, folder := range f.Subfolders {
		if found := folder.GetSubfolder(folderPath); found != nil {
			return found
		}
	}
	return nil
}

// LockByPath locks a folder or file at the given path. Locking a folder locks all descendants.
func (f *Folder) LockByPath(path string) {
	// If this folder matches, lock entire subtree
	if f.Path == path {
		f.lockRecursive()
		return
	}
	// Otherwise, delegate to subfolders and files
	for _, sf := range f.Subfolders {
		sf.LockByPath(path)
	}
	for _, file := range f.Files {
		file.LockByPath(path)
	}
}

// lockRecursive locks this folder and all nested folders and files
func (f *Folder) lockRecursive() {
	f.Locked = true
	fmt.Printf("Folder '%s' locked\n", f.Path)
	for _, sf := range f.Subfolders {
		sf.lockRecursive()
	}
	for _, file := range f.Files {
		file.Locked = true
		fmt.Printf("File '%s' locked\n", file.Path)
	}
}

// UnlockByPath unlocks a folder or file at the given path. Unlocking a folder unlocks all descendants.
func (f *Folder) UnlockByPath(path string) {
	if f.Path == path {
		f.unlockRecursive()
		return
	}
	for _, sf := range f.Subfolders {
		sf.UnlockByPath(path)
	}
	for _, file := range f.Files {
		file.UnlockByPath(path)
	}
}

// unlockRecursive unlocks this folder and all nested folders and files
func (f *Folder) unlockRecursive() {
	f.Locked = false
	fmt.Printf("Folder '%s' unlocked\n", f.Path)
	for _, sf := range f.Subfolders {
		sf.unlockRecursive()
	}
	for _, file := range f.Files {
		file.Locked = false
		fmt.Printf("File '%s' unlocked\n", file.Path)
	}
}

// AddTagToFile tags a file in this folder or its subfolders
func (f *Folder) AddTagToFile(filePath, tagName string) bool {
	file := f.GetFile(filePath)
	if file != nil {
		file.Tags = append(file.Tags, tagName)
		return true
	}
	return false
}

// AddTagToSelf adds a tag to the folder itself
func (f *Folder) AddTagToSelf(tagID, tagName string) {
	f.Tags = append(f.Tags, tagName)
}

// RemoveTag removes a tag from this folder
func (f *Folder) RemoveTag(tag string) bool {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags = append(f.Tags[:i], f.Tags[i+1:]...)
			return true
		}
	}
	return false
}

// Display prints the folder tree with lock status and tags
func (f *Folder) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFolder: %s, Locked=%t\n", prefix, f.Name, f.Locked)
	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s\n", prefix, tag)
		}
	}
	for _, sf := range f.Subfolders {
		sf.Display(indent + 1)
	}
	for _, file := range f.Files {
		file.Display(indent + 1)
	}
}

// -------------------- File Methods --------------------
// LockByPath locks this file if the path matches
func (f *File) LockByPath(path string) {
	if f.Path == path {
		f.Locked = true
		fmt.Printf("File '%s' locked\n", f.Path)
	}
}

// UnlockByPath unlocks this file if the path matches
func (f *File) UnlockByPath(path string) {
	if f.Path == path {
		f.Locked = false
		fmt.Printf("File '%s' unlocked\n", f.Path)
	}
}

// RemoveTag removes a tag from this file
func (f *File) RemoveTag(tag string) bool {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags = append(f.Tags[:i], f.Tags[i+1:]...)
			return true
		}
	}
	return false
}

// Display prints file info with lock status, metadata, and tags
func (f *File) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFile: %s, Locked=%t\n", prefix, f.Name, f.Locked)
	if len(f.Metadata) > 0 {
		fmt.Printf("%s  Metadata:\n", prefix)
		for _, entry := range f.Metadata {
			fmt.Printf("%s    - %s: %s\n", prefix, entry.Key, entry.Value)
		}
	}
	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s\n", prefix, tag)
		}
	}
}
