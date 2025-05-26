package filesystem

import (
	"net/http"
)

func getCompositeHandler(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("id")
	managerName := r.URL.Query().Get("name")
	filePath := r.URL.Query().Get("path")

	composite := ConvertToComposite(managerID, managerName, filePath)
	if composite == nil {
		http.Error(w, "Could not build composite", http.StatusInternalServerError)
		return
	}
	composite.Display(0)
	// json.NewEncoder(w).Encode(composite)
}

func HandleRequests() {
	http.HandleFunc("/composite", getCompositeHandler)
	// ...other endpoints...
	http.ListenAndServe(":51000", nil)
}
