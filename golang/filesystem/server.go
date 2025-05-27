package filesystem

import (
	"net/http"
)

var handleComposite func(*Folder, int, string)

func getCompositeHandler(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("id")
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	composite := ConvertToComposite(managerID, managerName, filePath)
	if composite == nil {
		w.Write([]byte("false"))
		return
	}
	if handleComposite != nil {
		handleComposite(composite, 0, "")
	}
	w.Write([]byte("true")) //temporary return
}
func removeCompositeHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if handleComposite != nil {
		handleComposite(nil, 1, filePath)
	}
	w.Write([]byte("true")) //temporary return
}
func HandleRequests(storeFunc func(*Folder, int, string)) {
	handleComposite = storeFunc
	http.HandleFunc("/addDirectory", getCompositeHandler)
	http.HandleFunc("/removeDirectory", removeCompositeHandler)
	http.ListenAndServe(":51000", nil)
}
