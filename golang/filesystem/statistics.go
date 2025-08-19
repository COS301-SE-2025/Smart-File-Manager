package filesystem

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type file struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
}

type ManagerStatistics struct {
	ManagerName    string `json:"manager_name"`
	Size           int64  `json:"size"`
	Folders        int    `json:"folders"`
	Files          int    `json:"files"`
	Recent         []file `json:"recent"`
	Largest        []file `json:"largest"`
	Oldest         []file `json:"oldest"`
	UmbrellaCounts []int  `json:"umbrella_counts"`
}

type fileInfo struct {
	path     string
	name     string
	size     int64
	modTime  time.Time
	umbrella string
}

func StatHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("StatHandler: Starting statistics collection")
	mu.Lock() // grab the lock
	composites := make([]*Folder, len(Composites))
	copy(composites, Composites)
	mu.Unlock() // release immediately

	var managers []ManagerStatistics
	// mu.Lock()
	// defer mu.Unlock()
	if len(composites) == 0 {
		json.NewEncoder(w).Encode(struct{}{})
		return
	}
	for _, folder := range composites {

		manager := ManagerStatistics{
			ManagerName: folder.Name,
		}

		// Collect all file information for this manager
		allFiles := collectManagerFiles(folder)

		// Calculate statistics
		manager.Files = len(allFiles)
		manager.Folders = countFolders(folder)
		manager.Size = calculateTotalSize(allFiles)
		// log.Printf("\n")
		manager.UmbrellaCounts = calculateUmbrellaCounts(folder, folder.Name)

		// Get file rankings
		manager.Recent = getNewestFiles(allFiles, 5)
		manager.Oldest = getOldestFiles(allFiles, 5)
		manager.Largest = getLargestFiles(allFiles, 5)

		managers = append(managers, manager)
		log.Printf("StatHandler: Completed manager: %s", folder.Name)
	}

	log.Println("StatHandler: Encoding response to JSON")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(managers)
	log.Println("StatHandler: Request completed successfully")
}

func collectManagerFiles(folder *Folder) []fileInfo {
	var files []fileInfo

	collectFilesRecursive(folder, &files)

	return files
}

func collectFilesRecursive(folder *Folder, files *[]fileInfo) {
	for _, file := range folder.Files {
		// Add timeout check periodically
		select {
		case <-time.After(50 * time.Millisecond):
			log.Printf("Timeout warning: Processing file %s is taking longer than expected", file.Path)
		default:
		}

		info, err := os.Stat(file.Path)
		if err != nil {
			continue
		}

		umbrella := "Unknown"
		if ObjectMap != nil && ObjectMap[folder.Name] != nil {
			if objInfo, exists := ObjectMap[folder.Name][file.Path]; exists {
				umbrella = objInfo.umbrellaType
			}
		}

		*files = append(*files, fileInfo{
			path:     file.Path,
			name:     filepath.Base(file.Path),
			size:     info.Size(),
			modTime:  info.ModTime(),
			umbrella: umbrella,
		})
	}

	for _, subFolder := range folder.Subfolders {
		collectFilesRecursive(subFolder, files)
	}
}

func countFolders(folder *Folder) int {
	count := len(folder.Subfolders)
	for _, subFolder := range folder.Subfolders {
		count += countFolders(subFolder)
	}
	return count
}

func calculateTotalSize(files []fileInfo) int64 {
	var total int64
	for _, file := range files {
		total += file.size
	}
	return total
}

func calculateUmbrellaCounts(folder *Folder, managerName string) []int {
	// [Documents, Images, Music, Presentations, Videos, Spreadsheets, Archives, Unknown]
	counts := make([]int, 8)
	umbrellaOrder := []string{
		"Documents", "Images", "Music", "Presentations", "Videos", "Spreadsheets", "Archives", "Unknown",
	}

	// Ensure types are loaded
	LoadTypes(folder, managerName)

	umbrellaMap := ObjectMap[managerName]
	for _, obj := range umbrellaMap {
		found := false
		for i, umbrella := range umbrellaOrder {
			if obj.umbrellaType == umbrella {
				counts[i]++
				found = true
				break
			}
		}
		if !found {
			counts[7]++ // Unknown
		}
	}
	return counts
}

func getNewestFiles(files []fileInfo, limit int) []file {
	if len(files) == 0 {
		return []file{}
	}

	// Sort by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})

	result := make([]file, 0, limit)
	count := limit
	if count > len(files) {
		count = len(files)
	}

	for i := 0; i < count; i++ {
		result = append(result, file{
			FilePath: files[i].path,
			FileName: files[i].name,
		})
	}

	return result
}

func getOldestFiles(files []fileInfo, limit int) []file {
	if len(files) == 0 {
		return []file{}
	}

	// Sort by modification time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	result := make([]file, 0, limit)
	count := limit
	if count > len(files) {
		count = len(files)
	}

	for i := 0; i < count; i++ {
		result = append(result, file{
			FilePath: files[i].path,
			FileName: files[i].name,
		})
	}

	return result
}

func getLargestFiles(files []fileInfo, limit int) []file {
	if len(files) == 0 {
		return []file{}
	}

	// Sort by size (largest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	result := make([]file, 0, limit)
	count := limit
	if count > len(files) {
		count = len(files)
	}

	for i := 0; i < count; i++ {
		result = append(result, file{
			FilePath: files[i].path,
			FileName: files[i].name,
		})
	}

	return result
}
