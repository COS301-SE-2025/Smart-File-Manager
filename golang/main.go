package main

import (
	"github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

func main() {

	//testing directory creation
	// folder := mockFolderStructure()
	// filesystem.CreateDirectoryStructure(folder)
	filesystem.HandleRequests()

	// print current Composites
	Composites := filesystem.GetComposites()
	for _, item := range Composites {
		item.Display(0)
	}
	// print current composites
	// composites := filesystem.GetComposites()
	// for _, item := range composites {
	// 	item.Display(0)
	// }
}
func mockFolderStructure() *filesystem.Folder {
	return &filesystem.Folder{
		NewPath: "test_root",
		Subfolders: []*filesystem.Folder{
			{
				NewPath: "test_root/sub1",
				Subfolders: []*filesystem.Folder{
					{
						NewPath: "test_root/sub1/sub1_1",
					},
				},
			},
			{
				NewPath: "test_root/sub2",
			},
		},
	}
}
