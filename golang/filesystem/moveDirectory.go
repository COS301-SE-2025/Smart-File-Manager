package filesystem

import (
	"fmt"
	"net/http"
	"os"
)

var root string = getPath()

func moveDirectoryHandler(w http.ResponseWriter, r *http.Request) {

	composites := GetComposites()
	if len(composites) == 0 {
		http.Error(w, "No managers found", http.StatusNotFound)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	compositeName := r.URL.Query().Get("name")

	for _, item := range composites {
		if item.Name == compositeName {
			//build the directory structure for the smart manager
			CreateDirectoryStructure(item)
			fmt.Println("Directory structure created for composite:", compositeName)
			//Move the content of the smart manager
			moveContent(item)
			fmt.Println("Content moved for composite:", compositeName)
			w.Write([]byte("true"))
			return
		}
	}
	fmt.Println("Smart manager not found: ", compositeName)
	w.Write([]byte("false"))

}

func moveContent(item *Folder) {
	//Move the files&folders according to new path after sorting
}

func CreateDirectoryStructure(item *Folder) {

}

func getPath() string { // navigate to the correct working directory

	err := os.Chdir("..")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return "NOPTH"
	}

	cwd, _ := os.Getwd()
	return cwd
}
