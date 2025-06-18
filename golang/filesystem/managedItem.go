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
}

// Folder structure
type Folder struct {
	ID           string
	Name         string
	Path         string
	newPath      string
	CreationDate time.Time
	Locked       bool
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
	fmt.Printf("%sFolder: %s\n", prefix, f.Name)

	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag)
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
func (f *File) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFile: %s\n", prefix, f.Name)

	if len(f.Metadata) > 0 {
		fmt.Printf("%s  Metadata:\n", prefix)
		for _, entry := range f.Metadata {
			fmt.Printf("%s    - %s: %s\n", prefix, entry.Key, entry.Value)
		}
	}

	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag)
		}
	}
}
