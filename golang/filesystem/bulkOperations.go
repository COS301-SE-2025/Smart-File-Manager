package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// /Type return structs
type object struct {
	fileType     string
	umbrellaType string
}

var FileTypeMap = map[string]string{
	// Documents
	"pdf": "Documents", "doc": "Documents", "docx": "Documents", "dot": "Documents",
	"dotx": "Documents", "rtf": "Documents", "txt": "Documents", "odt": "Documents",
	"ott": "Documents", "wpd": "Documents", "wps": "Documents", "md": "Documents",
	"log": "Documents", "tex": "Documents", "epub": "Documents", "mobi": "Documents",
	"azw": "Documents", "azw3": "Documents", "djvu": "Documents", "chm": "Documents",
	"ps": "Documents", "csv": "Documents",

	// Images
	"jpg": "Images", "jpeg": "Images", "png": "Images", "gif": "Images",
	"bmp": "Images", "tiff": "Images", "tif": "Images", "webp": "Images",
	"heic": "Images", "heif": "Images", "svg": "Images", "ico": "Images",
	"psd": "Images", "ai": "Images", "eps": "Images", "raw": "Images",
	"cr2": "Images", "nef": "Images", "orf": "Images", "arw": "Images",
	"dng": "Images",

	// Music
	"mp3": "Music", "wav": "Music", "flac": "Music", "aac": "Music",
	"ogg": "Music", "oga": "Music", "m4a": "Music", "wma": "Music",
	"opus": "Music", "aiff": "Music", "alac": "Music", "mid": "Music",
	"midi": "Music", "amr": "Music", "dsf": "Music", "dff": "Music",

	// Presentations
	"ppt": "Presentations", "pptx": "Presentations", "pps": "Presentations",
	"ppsx": "Presentations", "odp": "Presentations", "otp": "Presentations",
	"key": "Presentations",

	// Videos
	"mp4": "Videos", "m4v": "Videos", "mkv": "Videos", "avi": "Videos",
	"mov": "Videos", "wmv": "Videos", "flv": "Videos", "webm": "Videos",
	"vob": "Videos", "ts": "Videos", "m2ts": "Videos", "3gp": "Videos",
	"f4v": "Videos", "mpeg": "Videos", "mpg": "Videos", "ogv": "Videos",
	"divx": "Videos",

	// Spreadsheets
	"xls": "Spreadsheets", "xlsx": "Spreadsheets", "xlsm": "Spreadsheets",
	"ods": "Spreadsheets", "ots": "Spreadsheets", "tsv": "Spreadsheets",
	"xlsb": "Spreadsheets",

	// Archives
	"zip": "Archives", "rar": "Archives", "7z": "Archives", "tar": "Archives",
	"gz": "Archives", "bz2": "Archives", "xz": "Archives", "tgz": "Archives",
	"tbz2": "Archives", "lz": "Archives", "lzma": "Archives", "z": "Archives",
	"iso": "Archives", "dmg": "Archives", "cab": "Archives", "arj": "Archives",
	"ace": "Archives", "uue": "Archives",
}

type returnStruct struct {
	FilePath string   `json:"file_path"`
	FileName string   `json:"file_name"`
	FileTags []string `json:"file_tags,omitempty"`
}

var objectMap = map[string]map[string]object{}

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
			children := GoSidecreateDirectoryJSONStructure(folder)

			root := DirectoryTreeJson{
				Name:     folder.Name,
				IsFolder: true,
				RootPath: folder.Path,
				Children: children,
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
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
			children := GoSidecreateDirectoryJSONStructure(folder)

			root := DirectoryTreeJson{
				Name:     folder.Name,
				IsFolder: true,
				RootPath: folder.Path,
				Children: children,
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
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
			children := GoSidecreateDirectoryJSONStructure(folder)

			root := DirectoryTreeJson{
				Name:     folder.Name,
				IsFolder: true,
				RootPath: folder.Path,
				Children: children,
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
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
			children := GoSidecreateDirectoryJSONStructure(folder)

			root := DirectoryTreeJson{
				Name:     folder.Name,
				IsFolder: true,
				RootPath: folder.Path,
				Children: children,
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Encode the response as JSON
			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
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
			objectMap[name] = make(map[string]object)
			LoadTypes(c, name) // load types into the global objectMap
			var returnList []returnStruct
			// fmt.Println("LIST CREATED")
			if umbrella == "true" {
				for k, v := range objectMap[name] {
					if v.umbrellaType == types {
						returnList = append(returnList, returnStruct{
							FilePath: k,
							FileName: c.GetFile(k).Name,
							FileTags: c.GetFile(k).Tags,
						})
					}
				}
			} else {
				for k, v := range objectMap[name] {
					if v.fileType == types {
						// fmt.Println("FILE FOUND:", k, v.fileType)
						returnList = append(returnList, returnStruct{
							FilePath: k,
							FileName: c.GetFile(k).Name,
							FileTags: c.GetFile(k).Tags,
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

func LoadTypes(item *Folder, name string) {

	for _, file := range item.Files {
		objectMap[name][file.Path] = object{
			fileType:     strings.Split(file.Name, ".")[1],
			umbrellaType: GetCategory(file.Path),
		}
	}
	for _, subfolder := range item.Subfolders {
		LoadTypes(subfolder, name)
	}
}

func GetCategory(filename string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if category, exists := FileTypeMap[ext]; exists {
		return category
	}
	return "Unknown"
}
