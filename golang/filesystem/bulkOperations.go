package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// /Type return structs
type object struct {
	fileType     string
	umbrellaType string
}

type returnStruct struct {
	filePath string
	fileName string
}

var objectMap = map[string]object{}

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

// expected json for bulk delete
// [
//
//	{
//	  "file_path": "/home/user/documents,
//	},
//	{
//	  "file_path": "/home/user/photos,
//	},
//	{
//	  "file_path": "/home/user/music,
//	}
//
// ]
type DeleteStruct struct {
	FilePath string `json:"file_path"`
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
	for _, folder := range Composites {
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
	for _, folder := range Composites {
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

// DO NOT DELETE SUBFOLDERS ALONG WITH FOLDER THIS FUNCTION HANDLES IT
func BulkDeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var bulkList []DeleteStruct
	if err := json.NewDecoder(r.Body).Decode(&bulkList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var filePaths []string
	for _, item := range bulkList {
		filePaths = append(filePaths, item.FilePath)
		// fmt.Println("File path to delete:", item.FilePath)
	}

	mu.Lock()
	defer mu.Unlock()

	//delete all folders in list
	for _, folder := range Composites {
		if folder.Name == name {
			err := folder.RemoveMultipleSubfolders(filePaths)
			for _, path := range filePaths {
				if err[path] != nil {
					fmt.Println("Error removing file:", path, "Error:", err[path])
					continue
				}
				err := os.RemoveAll(path)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to remove folder %s: %v", path, err), http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Folders removed successfully"))
			return
		}
	}

	http.Error(w, "Folder not found", http.StatusNotFound)
}

func BulkDeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var bulkList []DeleteStruct
	if err := json.NewDecoder(r.Body).Decode(&bulkList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var filePaths []string
	for _, item := range bulkList {
		filePaths = append(filePaths, item.FilePath)
		// fmt.Println("File path to delete:", item.FilePath)
	}

	mu.Lock()
	defer mu.Unlock()

	//delete all folders in list
	for _, folder := range Composites {
		if folder.Name == name {
			err := folder.RemoveMultipleFiles(filePaths)

			for _, path := range filePaths {
				if err[path] != nil {
					fmt.Println("Error removing file:", path, "Error:", err[path])
					continue
				}
				err := os.RemoveAll(path)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to remove files %s: %v", path, err), http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Files removed successfully"))
			return
		}
	}

	http.Error(w, "Files not found", http.StatusNotFound)
}

func ReturnTypeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	types := r.URL.Query().Get("type")
	umbrella := r.URL.Query().Get("umbrella")
	w.Header().Set("Content-Type", "application/json")
	if name == "" || types == "" || umbrella == "" {
		http.Error(w, "Missing 'name' or 'type' or 'umbrella' parameter", http.StatusBadRequest)
		return
	}
	for _, c := range Composites {
		if c.Name == name {
			var returnList []returnStruct
			// fmt.Println("LIST CREATED")
			if umbrella == "true" {
				for k, v := range objectMap {
					if v.umbrellaType == types {
						returnList = append(returnList, returnStruct{
							filePath: k,
							fileName: c.GetFile(k).Name,
						})
					}
				}
			} else {
				for k, v := range objectMap {
					if v.fileType == types {
						// fmt.Println("FILE FOUND:", k, v.fileType)
						returnList = append(returnList, returnStruct{
							filePath: k,
							fileName: c.GetFile(k).Name,
						})
					}

				}
			}
			// w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(returnList)
			return
		}

	}
	http.Error(w, "Composite not found", http.StatusNotFound)

}

func LoadTypes(item *Folder) {
	for _, file := range item.Files {
		objectMap[file.Path] = object{
			fileType:     strings.Split(file.Name, ".")[1],
			umbrellaType: getUmbrellaType(file.Path),
		}
	}
	for _, subfolder := range item.Subfolders {
		LoadTypes(subfolder)
	}
}

func getUmbrellaType(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return ""
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return ""
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType
}
