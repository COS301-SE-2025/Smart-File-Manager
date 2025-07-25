package filesystem

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

const limit int = 25
const maxDist int = 25

func levelshteinDist(searchText string, fileName string) int {
	if fileName[0] != '.' {

		searchText = strings.ToLower(searchText)
		fileName = strings.ToLower(fileName)

		//if the search doesnt contain a . then we remove the file
		// extention from the fie names for better search results
		if !strings.Contains(searchText, ".") {

			if i := strings.LastIndex(searchText, "."); i >= 0 {
				searchText = searchText[:i]
			}
			if i := strings.LastIndex(fileName, "."); i >= 0 {
				fileName = fileName[:i]
			}
		}
	}

	// ── BOOST exact substrings ──
	if strings.Contains(fileName, searchText) {
		return 1
	}

	// now fall back on full Levenshtein
	la, lb := len(searchText), len(fileName)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// ensure la >= lb
	if la < lb {
		searchText, fileName = fileName, searchText
		la, lb = lb, la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		ai := searchText[i-1]
		for j := 1; j <= lb; j++ {
			cost := 0
			if ai != fileName[j-1] {
				cost = 1
			}
			sub := prev[j-1] + cost
			ins := curr[j-1] + 1
			del := prev[j] + 1

			// take the minimum
			if ins < sub {
				sub = ins
			}
			if del < sub {
				sub = del
			}
			curr[j] = sub
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

type safeResults struct {
	Name        string
	rankedFiles []rankedFile
}

type rankedFile struct {
	file     File
	distance int
}

type folderResponse struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	compositeName := r.URL.Query().Get("compositeName")
	searchText := r.URL.Query().Get("searchText")

	for _, comp := range composites {
		if comp.Name == compositeName {

			sr := getMatches(searchText, comp)

			// build a FolderResponse with just the File objects
			resp := folderResponse{
				Name:  sr.Name,
				Files: make([]File, len(sr.rankedFiles)),
			}
			for i, rf := range sr.rankedFiles {
				resp.Files[i] = rf.file
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

}

func getMatches(text string, composite *Folder) *safeResults {
	res := &safeResults{
		Name: composite.Name,
	}

	resultChan := make(chan rankedFile)

	var wg sync.WaitGroup

	wg.Add(1)
	go exploreFolder(composite, text, resultChan, &wg)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	seen := make(map[string]struct{})
	for currentRankedFile := range resultChan {

		key := currentRankedFile.file.Path
		if _, ok := seen[key]; ok {
			continue // already inserted this file
		}
		seen[key] = struct{}{}

		inserted := false
		for i, iteratedRankedFile := range res.rankedFiles {
			//if current is better
			if iteratedRankedFile.distance > currentRankedFile.distance {

				if len(res.rankedFiles) < limit {
					//insert by shifting over to make the array sorted
					res.rankedFiles = append(
						append(res.rankedFiles[:i], currentRankedFile),
						res.rankedFiles[i:]...,
					)

					inserted = true
					break
				} else {
					res.rankedFiles = append(
						append(
							res.rankedFiles[:i],
							currentRankedFile,
						),
						res.rankedFiles[i:limit-1]...,
					)
					inserted = true
					break
				}
			}
		}
		// if limit is not reached and not inserted then we insert at the end
		if len(res.rankedFiles) < limit {
			if !inserted {
				res.rankedFiles = append(res.rankedFiles, currentRankedFile)
			}
		}
	}
	unique := make([]rankedFile, 0, len(res.rankedFiles))
	finalSeen := make(map[string]struct{}, len(res.rankedFiles))
	for _, rf := range res.rankedFiles {
		key := filepath.Clean(rf.file.Path)
		if _, ok := finalSeen[key]; ok {
			continue
		}
		finalSeen[key] = struct{}{}
		unique = append(unique, rf)
	}
	res.rankedFiles = unique

	return res
}

func exploreFolder(f *Folder, text string, c chan<- rankedFile, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, folder := range f.Subfolders {
		wg.Add(1)
		go exploreFolder(folder, text, c, wg)
	}
	for _, file := range f.Files {
		dist := levelshteinDist(text, file.Name)
		if dist <= maxDist {

			c <- rankedFile{file: *file, distance: dist}
		}
	}
}
