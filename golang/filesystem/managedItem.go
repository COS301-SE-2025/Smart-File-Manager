package filesystem

import (
	"fmt"
	"strings"
	"time"
)

// Tag structure
type Tag struct {
	ID   string
	Name string
}

// File structure
type File struct {
	Name     string
	Path     string
	Metadata map[string]string
	Tags     []Tag
}

// Folder structure
type Folder struct {
	ID           string
	Name         string
	Path         string
	CreationDate time.Time
	Locked       bool
	Files        []*File
	Subfolders   []*Folder
	Tags         []Tag
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

// AddTagToFile adds a tag to a file inside the folder
func (f *Folder) AddTagToFile(filePath, tagID, tagName string) bool {
	file := f.GetFile(filePath)
	if file != nil {
		file.Tags = append(file.Tags, Tag{ID: tagID, Name: tagName})
		return true
	}
	return false
}

// AddTagToSelf adds a tag directly to this folder
func (f *Folder) AddTagToSelf(tagID, tagName string) {
	f.Tags = append(f.Tags, Tag{ID: tagID, Name: tagName})
}

// Display the folder and its contents recursively
func (f *Folder) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFolder: %s\n", prefix, f.Name)

	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag.ID, tag.Name)
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

// Display file and its metadata
func (f *File) Display(indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sFile: %s\n", prefix, f.Name)

	if len(f.Metadata) > 0 {
		fmt.Printf("%s  Metadata:\n", prefix)
		for key, value := range f.Metadata {
			fmt.Printf("%s    - %s: %s\n", prefix, key, value)
		}
	}

	if len(f.Tags) > 0 {
		fmt.Printf("%s  Tags:\n", prefix)
		for _, tag := range f.Tags {
			fmt.Printf("%s    - %s: %s\n", prefix, tag.ID, tag.Name)
		}
	}
}
