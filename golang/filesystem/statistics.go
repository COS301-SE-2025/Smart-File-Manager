package filesystem

import (
	"log"
	"os"
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
			"umbrella_percentage": ["5%","20%","10%","20%","10%","20%","10%","15%"]
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
			"umbrella_percentage": ["5%","20%","10%","20%","10%","20%","10%","15%"]
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
			"umbrella_percentage": ["5%","20%","10%","20%","10%","20%","10%","15%"]
	//	}
	//

// ]
*/
type file struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
}
type ManagerStatistics struct {
	ManagerName        string   `json:"manager_name"`
	Size               int64    `json:"size"`
	Folders            int      `json:"folders"`
	Files              int      `json:"files"`
	Recent             []file   `json:"recent"`
	Largest            []file   `json:"largest"`
	Oldest             []file   `json:"oldest"`
	UmbrellaPercentage []string `json:"umbrella_percentage"`
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
