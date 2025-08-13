package filesystem

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

/*
	statistics we need
	num managers
	num files	& per manager
	num folders & per manager
	size of manager & per manager

	Most recently used files	lim=5
	largest files	lim=5
	least recently used files 	lim=5


	umbrella files perentage of all files

	json structure:
	// [
	//	{
	//	 	"manager_name": manager1,
			"size": 1024,
			"folders": 10,
			"files": 100,
			"recent": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"largest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"oldest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"umbrella_percentage": ["5","20","10","20","10","20","10","15"]
	//	},
	//	{
	//	  	"manager_name": manager2,
			"size": 1024,
			"folders": 10,
			"files": 100,
			"recent": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"largest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"oldest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"umbrella_percentage": ["5","20","10","20","10","20","10","15"]
	//	},
	//	{
	//	  	"manager_name": manager3,
			"size": 1024,
			"folders": 10,
			"files": 100,
			"recent": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"largest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"oldest": {
						file_paths: ["/home/user/documents/report.pdf", "/home/user/photos/vacation.jpg", "/home/user/music/song.mp3"],
						file_names: ["report.pdf", "vacation.jpg", "song.mp3"]
						},
			"umbrella_percentage": ["5","20","10","20","10","20","10","15"]
	//	}
	//

// ]
*/
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

func StatHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	//load all maps
	assignMaps()
	//create object that will be converted into json
	managers := make([]ManagerStatistics, 0)
	for _, folder := range GetComposites() {
		manager := ManagerStatistics{
			ManagerName: folder.Name,
		}
		getNumItems(&manager, folder)
		getManagerSize(&manager, folder)
		getUmbrellaRatio(&manager, folder)
		getFileStats(&manager)
		managers = append(managers, manager)
	}
	//convert to json
	json.NewEncoder(w).Encode(managers)
	return

}

func assignMaps() {
	//initialize if needed
	if len(fileStatsReturn.Newest) == 0 {
		fileStatsReturn.Newest = []file{}
	}
	if len(fileStatsReturn.Oldest) == 0 {
		fileStatsReturn.Oldest = []file{}
	}
	if len(fileStatsReturn.Largest) == 0 {
		fileStatsReturn.Largest = []file{}
	}
	for _, folder := range GetComposites() {
		readIntoMap(&timesMap, &sizeMap, folder)
		LoadTypes(folder, folder.Name)
	}
	assignTimedFiles()
	assignSizeFiles()

}

func getNumItems(m *ManagerStatistics, item *Folder) {
	m.Files += len(item.Files)
	for _, subFolder := range item.Subfolders {
		getNumItems(m, subFolder)
	}
}

func getManagerSize(m *ManagerStatistics, item *Folder) {
	for _, file := range item.Files {
		// Get file information using os.Stat()
		fileInfo, err := os.Stat(file.Path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("Error: File '%s' does not exist.\n", file.Path)
			} else {
				log.Fatalf("Error getting file info for '%s': %v\n", file.Path, err)
			}
			return
		}

		// Access the file size from the FileInfo
		fileSize := fileInfo.Size()
		m.Size += fileSize
	}

	for _, subFolder := range item.Subfolders {
		getManagerSize(m, subFolder)
	}
}

// Documents
// Images
// Music
// Presentations
// Videos
// Spreadsheets
// Archives
// Unknown
// remember to call loadTypes before calling this function
func getUmbrellaRatio(m *ManagerStatistics, item *Folder) {
	umbrellaCounts := make([]int, 8)
	for _, file := range item.Files {
		umbrellaType := ObjectMap[item.Name][file.Path].umbrellaType
		switch umbrellaType {
		case "Documents":
			umbrellaCounts[0]++
		case "Images":
			umbrellaCounts[1]++
		case "Music":
			umbrellaCounts[2]++
		case "Presentations":
			umbrellaCounts[3]++
		case "Videos":
			umbrellaCounts[4]++
		case "Spreadsheets":
			umbrellaCounts[5]++
		case "Archives":
			umbrellaCounts[6]++
		default:
			umbrellaCounts[7]++
		}
	}

	m.UmbrellaCounts = umbrellaCounts
}

type fileStats struct {
	Newest  []file `json:"newest"`
	Oldest  []file `json:"oldest"`
	Largest []file `json:"largest"`
}

var fileStatsReturn fileStats

func getFileStats(m *ManagerStatistics) {
	m.Recent = fileStatsReturn.Newest
	m.Oldest = fileStatsReturn.Oldest
	m.Largest = fileStatsReturn.Largest
}

var timesMap map[float64]string
var sizeMap map[float64]string

func readIntoMap(times *map[float64]string, sizes *map[float64]string, f *Folder) {
	for _, item := range f.Files {
		// Get the file modification time
		fileInfo, err := os.Stat(item.Path)
		if err != nil {
			log.Printf("Error getting file info for '%s': %v\n", item.Path, err)
			continue
		}
		modTime := fileInfo.ModTime().UnixNano()
		timesMap[float64(modTime)] = item.Path
		fileSize := fileInfo.Size()
		sizeMap[float64(fileSize)] = item.Path
	}
	for _, subFolder := range f.Subfolders {
		readIntoMap(times, sizes, subFolder)
	}
}

func assignTimedFiles() {
	basket := timesMap
	keys := make([]float64, 0, len(basket))
	for k := range basket {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for i := 0; i < 5; i++ {
		filePath := basket[keys[i]]
		fileStatsReturn.Newest = append(fileStatsReturn.Newest, file{FilePath: filePath, FileName: filepath.Base(filePath)})
	}
	for i := len(keys) - 6; i > len(keys)-1; i-- {
		filePath := basket[keys[i]]
		fileStatsReturn.Oldest = append(fileStatsReturn.Oldest, file{FilePath: filePath, FileName: filepath.Base(filePath)})
	}
}

func assignSizeFiles() {
	basket := sizeMap
	keys := make([]float64, 0, len(basket))
	for k := range basket {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for i := 0; i < 5; i++ {
		filePath := basket[keys[i]]
		fileStatsReturn.Largest = append(fileStatsReturn.Largest, file{FilePath: filePath, FileName: filepath.Base(filePath)})
	}
}
