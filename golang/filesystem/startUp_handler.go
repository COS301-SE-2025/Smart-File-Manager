package filesystem

// calling /startUp will load to memory the managers that have been created already by the user
// this means you wont have to create a new one each time you open the app. no parameters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type ManagerRecord struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type startUpResponse struct {
	ResponseMessage string   `json:"responseMessage"`
	ManagerNames    []string `json:"managerNames"`
}

var managersFilePath = filepath.Join("storage", "startUpStorageFile.json")

// used to change directory during testing
func SetManagersFilePath(p string) {
	managersFilePath = p
}

// api entry
func startUpHandler(w http.ResponseWriter, r *http.Request) {

	Composites = nil

	recs, err := loadManagerRecords()

	if err != nil {
		errMsg := "Internal error: " + err.Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Clean up any JSON files in storage/ that aren't for a known manager.
	if err := cleanupOrphanCompositeJSONs(recs); err != nil {
		// Best-effort cleanup; don't fail startup
		fmt.Printf("cleanupOrphanCompositeJSONs warning: %v\n", err)
	}

	var (
		managerNames []string
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	for _, r := range recs {
		wg.Add(1)
		go func(rec ManagerRecord) {
			defer wg.Done()

			composite, err := ConvertToObject(rec.Name, rec.Path)

			if err != nil {
				fmt.Printf("ConvertToObject failed for %s (%s) in %v\n", rec.Name, rec.Path, err)
				return
			}

			mu.Lock()
			Composites = append(Composites, composite)

			managerNames = append(managerNames, composite.Name)
			mu.Unlock()
		}(r)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)

	res := startUpResponse{
		ResponseMessage: "Request successful!, Composites: " + strconv.Itoa(len(managerNames)),
		ManagerNames:    managerNames,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func loadManagerRecords() ([]ManagerRecord, error) {
	// If the file doesn't exist yet, start with empty
	data, err := os.ReadFile(managersFilePath)

	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var recs []ManagerRecord

	//populates recs
	if err := json.Unmarshal(data, &recs); err != nil {
		fmt.Println("error in unmarshaling of json")
		//delete save file to restart managers
		os.Remove(managersFilePath)
		return nil, err
	}
	return recs, nil
}

// writes to the json file that tracks which managers exist
func saveManagerRecords(recs []ManagerRecord) error {
	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(managersFilePath), 0755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(managersFilePath, out, 0644)
}

// functions used when adding/removing managers that keeps track of the ones to save:

func AddManager(name, path string) error {
	composite, err := ConvertToObject(name, path)

	saveCompositeDetails(composite)

	// fmt.Println("+++++++++++++")
	// PrettyPrintFolder(composite, "")
	// fmt.Println("+++++++++++++")
	if err != nil {
		return err
	}
	Composites = append(Composites, composite)

	// rebuild the small record slice
	var recs []ManagerRecord
	for _, f := range Composites {
		recs = append(recs, ManagerRecord{Name: f.Name, Path: f.Path})
	}

	return saveManagerRecords(recs)
}

func RemoveManager(path string) error {
	// remove from Composites by Path, then rewrite file
	var recs []ManagerRecord
	for i, f := range Composites {
		if f.Path == path {
			Composites = append(Composites[:i], Composites[i+1:]...)
			break
		}
	}
	for _, f := range Composites {
		recs = append(recs, ManagerRecord{Name: f.Name, Path: f.Path})
	}
	if err := saveManagerRecords(recs); err != nil {
		return err
	}

	if err := cleanupOrphanCompositeJSONs(recs); err != nil {
		fmt.Printf("cleanupOrphanCompositeJSONs warning: %v\n", err)
	}
	return nil

}

func cleanupOrphanCompositeJSONs(recs []ManagerRecord) error {
	dir := filepath.Dir(managersFilePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Build allowed set: manager JSONs + the managers file itself.
	allowed := make(map[string]struct{}, len(recs)+1)
	allowed[filepath.Base(managersFilePath)] = struct{}{}
	for _, r := range recs {
		if r.Name != "" {
			allowed[r.Name+".json"] = struct{}{}
		}
	}

	var firstErr error
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}
		if _, ok := allowed[name]; ok {
			continue
		}
		path := filepath.Join(dir, name)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}
