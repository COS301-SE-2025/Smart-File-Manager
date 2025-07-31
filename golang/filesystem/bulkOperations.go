package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// this file will contain bulk operations for the filesystem such as bulk add, delete of files, folders and adding tags
// expected json for bulk add tags
// [
//
//	{
//	  "file_path": "/home/user/documents/report.pdf",
//	  "tags": ["work", "important", "pdf"]
//	},
//	{
//	  "file_path": "/home/user/photos/vacation.jpg",
//	  "tags": ["holiday", "family", "2025"]
//	},
//	{
//	  "file_path": "/home/user/music/song.mp3",
//	  "tags": ["music", "mp3", "favorites"]
//	}
//
// ]
// struct for tag json
type TagsStruct struct {
	FilePath string   `json:"file_path"`
	Tags     []string `json:"tags"`
}

// BulkAddTags adds tags to files in bulk
func BulkAddTags(item *Folder, bulkList []TagsStruct) error {
	for _, tagItem := range bulkList {
		for _, tag := range tagItem.Tags {
			item.AddTagToFile(tagItem.FilePath, tag)
		}
	}
	return nil
}

// BulkRemoveTags removes tags from files in bulk
func BulkRemoveTags(item *Folder, bulkList []TagsStruct) error {
	for _, tagItem := range bulkList {
		file := item.GetFile(tagItem.FilePath)
		if file == nil {
			return fmt.Errorf("file not found: %s", tagItem.FilePath)
		}
		for _, tag := range tagItem.Tags {
			if !file.RemoveTag(tag) {
				return fmt.Errorf("failed to remove tag %s from file %s", tag, tagItem.FilePath)
			}
		}
	}
	return nil
}

func BulkAddTagHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var bulkList []TagsStruct
	if err := json.NewDecoder(r.Body).Decode(&bulkList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find the corresponding Folder by name
	for _, folder := range composites {
		if folder.Name == name {
			if err := BulkAddTags(folder, bulkList); err != nil {
				http.Error(w, fmt.Sprintf("Failed to add tags: %v", err), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Tags added successfully"))
			return
		}
	}

	http.Error(w, "Folder not found", http.StatusNotFound)
}

func BulkRemoveTagHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var bulkList []TagsStruct
	if err := json.NewDecoder(r.Body).Decode(&bulkList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find the corresponding Folder by name
	for _, folder := range composites {
		if folder.Name == name {
			if err := BulkRemoveTags(folder, bulkList); err != nil {
				http.Error(w, fmt.Sprintf("Failed to remove tags: %v", err), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Tags removed successfully"))
			return
		}
	}

	http.Error(w, "Folder not found", http.StatusNotFound)
}
