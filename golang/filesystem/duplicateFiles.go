package filesystem

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func FindDuplicateFiles(root *Folder) []DuplicateEntry {
	// 1st pass: group all files by their file size
	sizeBuckets := make(map[int64][]string)
	collectBySize(root, sizeBuckets)

	duplicates := []DuplicateEntry{}

	// 2nd pass: for each bucket with >1 file, compute full hash and detect dups
	for _, paths := range sizeBuckets {
		if len(paths) < 2 {
			continue
		}

		hashMap := make(map[string]string)
		for _, p := range paths {
			hash := computeFileHash(p)
			if hash == "" {
				continue
			}
			if orig, exists := hashMap[hash]; exists {
				// duplicate found
				duplicates = append(duplicates, DuplicateEntry{
					Name:      filepath.Base(p),
					Original:  orig,
					Duplicate: p,
				})
			} else {
				hashMap[hash] = p
			}
		}
	}

	return duplicates
}

// collectBySize recurses through folder tree, grouping file paths by size
func collectBySize(folder *Folder, buckets map[int64][]string) {
	for _, f := range folder.Files {
		if info, err := os.Stat(f.Path); err == nil && info.Mode().IsRegular() {
			buckets[info.Size()] = append(buckets[info.Size()], f.Path)
		}
	}
	for _, sub := range folder.Subfolders {
		collectBySize(sub, buckets)
	}
}

func findDuplicateFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	name := r.URL.Query().Get("name")
	mu.Lock()
	defer mu.Unlock()

	for _, c := range Composites {
		if c.Name == name {
			duplicates := FindDuplicateFiles(c)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(duplicates)
			return
		}
	}

	http.Error(w, "composite not found", http.StatusNotFound)
}
