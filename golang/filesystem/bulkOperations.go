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

type Taggable interface {
	AddTagToFile(path string, tag string) bool
}

func BulkAddTags(item Taggable, bulkList []TagsStruct) error {
	for _, tagItem := range bulkList {
		for _, tag := range tagItem.Tags {
			item.AddTagToFile(tagItem.FilePath, tag)
		}
	}
	return nil
}

func BulkTagHandler(w http.ResponseWriter, r *http.Request) {
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
