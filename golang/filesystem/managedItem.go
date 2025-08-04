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
	NewPath  string
	Metadata []*MetadataEntry
	Tags     []string
	Locked   bool // Lock status for file
}

// Folder represents a directory in the filesystem
type Folder struct {
	Name         string
	Path         string
	NewPath      string
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

func (f *Folder) RemoveFile(filePath string) error {
	if f.Locked {
		return fmt.Errorf("cannot remove file: folder '%s' is locked", f.Name)
	}

	for i, file := range f.Files {
		if file.Path == filePath {
			f.Files[i] = f.Files[len(f.Files)-1]
			f.Files[len(f.Files)-1] = nil
			f.Files = f.Files[:len(f.Files)-1]
			return nil
		}
	}

	for _, subfolder := range f.Subfolders {
		if err := subfolder.RemoveFile(filePath); err == nil {
			return nil
		}
	}

	return fmt.Errorf("file not found: %s", filePath)
}

func (f *Folder) RemoveFileOrderPreserving(filePath string) error {

	for i, file := range f.Files {
		if file.Path == filePath {

			copy(f.Files[i:], f.Files[i+1:])
			f.Files[len(f.Files)-1] = nil
			f.Files = f.Files[:len(f.Files)-1]
			return nil
		}
	}

	for _, subfolder := range f.Subfolders {
		if err := subfolder.RemoveFileOrderPreserving(filePath); err == nil {
			return nil
		}
	}

	return fmt.Errorf("file not found: %s", filePath)
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
	checkFile := f.GetFile(path)
	if checkFile != nil {
		checkFile.Lock()
		return
	}
	checkFolder := f.GetSubfolder(path)
	if checkFolder != nil {
		checkFolder.lockRecursive()
		return
	}
}

// lockRecursive locks this folder and all nested folders and files
func (f *Folder) lockRecursive() {
	f.Locked = true
	for _, sf := range f.Subfolders {
		sf.lockRecursive()
	}
	for _, file := range f.Files {
		file.Locked = true
	}
}

// UnlockByPath unlocks a folder or file at the given path. Unlocking a folder unlocks all descendants.
func (f *Folder) UnlockByPath(path string) {
	checkFile := f.GetFile(path)
	if checkFile != nil {
		checkFile.Unlock()
		return
	}
	checkFolder := f.GetSubfolder(path)
	if checkFolder != nil {
		checkFolder.unlockRecursive()
		return
	}
}

// unlockRecursive unlocks this folder and all nested folders and files
func (f *Folder) unlockRecursive() {
	f.Locked = false
	for _, sf := range f.Subfolders {
		sf.unlockRecursive()
	}
	for _, file := range f.Files {
		file.Locked = false
	}
}

// AddTagToFile tags a file in this folder or its subfolders
func (f *Folder) AddTagToFile(filePath, tagName string) bool {
	file := f.GetFile(filePath)
	if file != nil {
		for _, tag := range file.Tags {
			if tag == tagName {
				return false // Tag already exists
			}
		}
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
// Lock locks this file
func (f *File) Lock() {
	f.Locked = true
}

// Unlock unlocks this file
func (f *File) Unlock() {
	f.Locked = false
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
