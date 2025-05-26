package filesystem

import (
	"net/http"
)

var storeCompositeFunc func(*Folder)

func getCompositeHandler(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("id")
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	composite := ConvertToComposite(managerID, managerName, filePath)
	if composite == nil {
		w.Write([]byte("false"))
		return
	}
	if storeCompositeFunc != nil {
		storeCompositeFunc(composite)
	}
	w.Write([]byte("true")) //temporary return
}

func HandleRequests(storeFunc func(*Folder)) {
	storeCompositeFunc = storeFunc
	http.HandleFunc("/composite", getCompositeHandler)
	http.ListenAndServe(":51000", nil)
}
