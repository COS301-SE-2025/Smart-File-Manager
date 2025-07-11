package filesystem

import (
	"fmt"
	"strings"
	"time"
)

type MetadataEntry struct {
	Key   string
	Value string
}

// File structure
type File struct {
	Name     string
	Path     string
	newPath  string
	Metadata []*MetadataEntry
	Tags     []string
	Locked   bool // Lock status for file
}

// Folder structure
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

// GetFile returns a file by path
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

// GetSubfolder returns a folder by path
func (f *Folder) GetSubfolder(folderPath string) *Folder {
	for _, folder := range f.Subfolders {
		if folder.Path == folderPath {
			return folder
		}
		if found := folder.GetSubfolder(folderPath); found != nil {
			return found
		}
	}
	return nil
}

// LockByPath locks a folder or file by its path
func (f *Folder) LockByPath(path string) {
	if f.Path == path {
		if !f.Locked {
			f.Locked = true
			fmt.Printf("Folder '%s' locked\n", f.Name)
		}
		for _, subfolder := range f.Subfolders {
			subfolder.LockByPath(path)
		}
		for _, file := range f.Files {
			file.LockByPath(path)
		}
	}
	for _, subfolder := range f.Subfolders {
		subfolder.LockByPath(path)
	}
}

// UnlockByPath unlocks a folder or file by its path
func (f *Folder) UnlockByPath(path string) {
	if f.Path == path {
		if f.Locked {
			f.Locked = false
			fmt.Printf("Folder '%s' unlocked\n", f.Name)
		}
		for _, subfolder := range f.Subfolders {
			subfolder.UnlockByPath(path)
		}
		for _, file := range f.Files {
			file.UnlockByPath(path)
		}
	}
	for _, subfolder := range f.Subfolders {
		subfolder.UnlockByPath(path)
	}
}

// Lock locks the file
func (f *File) LockByPath(path string) {
	if f.Path == path {
		if !f.Locked {
			f.Locked = true
			fmt.Printf("File '%s' locked\n", f.Name)
		}
	}
}

// Unlock unlocks the file
func (f *File) UnlockByPath(path string) {
	if f.Path == path {
		if f.Locked {
			f.Locked = false
			fmt.Printf("File '%s' unlocked\n", f.Name)
		}
	}
}

func (f *Folder) AddTagToFile(filePath, tagName string) bool {
	file := f.GetFile(filePath)
	if file != nil {
		file.Tags = append(file.Tags, tagName)
		return true
	}
	return false
}

func (f *Folder) AddTagToSelf(tagID, tagName string) {
	f.Tags = append(f.Tags, tagName)
}

func (f *Folder) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFolder: %s, Locked= %s\n", prefix, f.Name, fmt.Sprintf("%t", f.Locked))

	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s\n", prefix, tag)
		}
	}

	for _, sub := range f.Subfolders {
		sub.Display(indent + 1)
	}

	for _, file := range f.Files {
		file.Display(indent + 1)
	}
}

// -------------------- File Methods --------------------
func (f *File) RemoveTag(tag string) bool {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags = append(f.Tags[:i], f.Tags[i+1:]...)
			return true
		}
	}
	return false
}

func (f *Folder) RemoveTag(tag string) bool {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags = append(f.Tags[:i], f.Tags[i+1:]...)
			return true
		}
	}
	return false
}

// Display method for files
func (f *File) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFile: %s, Locked= %s\n", prefix, f.Name, fmt.Sprintf("%t", f.Locked))

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
