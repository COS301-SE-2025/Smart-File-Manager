package filesystem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type ManagerRecord struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const managersFile = "./storage/main.json"

// api entry
func startUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("setup called")
	recs, err := loadManagerRecords()

	if err != nil {
		var errMsg string = "Internal error: " + err.Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	//send the smart managers
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `{"message": "Request successful!, composites: %s"}`, strconv.Itoa(len(recs)))

	for _, r := range recs {
		// Convert the record back into your in‑memory Folder
		// using exactly the same logic as getCompositeHandler
		composite, err := ConvertToObject(r.Name, r.Path)
		if err != nil {
			fmt.Printf("readingComposite error: %s", err.Error())
			continue
		}
		composites = append(composites, composite)
	}

}

func loadManagerRecords() ([]ManagerRecord, error) {
	// If the file doesn't exist yet, start with empty
	data, err := os.ReadFile(managersFile)

	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var recs []ManagerRecord

	//populates recs
	if err := json.Unmarshal(data, &recs); err != nil {
		fmt.Println("error in unmarshaling of json")
		return nil, err
	}
	return recs, nil
}

// writes to the json file
func saveManagerRecords(recs []ManagerRecord) error {
	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(managersFile), 0755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(managersFile, out, 0644)
}

// functions used when adding/removing managers that keeps track of the ones to save:


func AddManager(name, path string) error {
	composite, err := ConvertToObject(name, path)
	if err != nil {
		return err
	}
	composites = append(composites, composite)

	// rebuild the small record slice
	var recs []ManagerRecord
	for _, f := range composites {
		recs = append(recs, ManagerRecord{Name: f.Name, Path: f.Path})
	}
	return saveManagerRecords(recs)
}

func RemoveManager(path string) error {
	// remove from composites by Path, then rewrite file
	var recs []ManagerRecord
	for i, f := range composites {
		if f.Path == path {
			composites = append(composites[:i], composites[i+1:]...)
			break
		}
	}
	for _, f := range composites {
		recs = append(recs, ManagerRecord{Name: f.Name, Path: f.Path})
	}
	return saveManagerRecords(recs)
}
