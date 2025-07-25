package filesystem

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type DuplicateEntry struct {
	Name      string `json:"name"`
	Original  string `json:"original"`
	Duplicate string `json:"duplicate"`
}

func computeFileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("failed to open file %s: %v", path, err)
		return ""
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		log.Printf("failed to hash file %s: %v", path, err)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func FindDuplicateFiles(item *Folder) []DuplicateEntry {
	fileHashes := make(map[string]string)
	var duplicates []DuplicateEntry

	var walk func(folder *Folder)
	walk = func(folder *Folder) {
		for _, file := range folder.Files {
			hash := computeFileHash(file.Path)
			if hash == "" {
				continue
			}
			if origPath, exists := fileHashes[hash]; exists {
				duplicates = append(duplicates, DuplicateEntry{
					Name:      file.Name,
					Original:  origPath,
					Duplicate: file.Path,
				})
			} else {
				fileHashes[hash] = file.Path
			}
		}
		for _, sub := range folder.Subfolders {
			walk(sub)
		}
	}

	walk(item)
	return duplicates
}

func findDuplicateFilesHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range composites {
		if c.Name == name {
			duplicates := FindDuplicateFiles(c)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(duplicates)
			return
		}
	}

	http.Error(w, "composite not found", http.StatusNotFound)
}
