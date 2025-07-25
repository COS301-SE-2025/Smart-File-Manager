package filesystem

import (
	"encoding/json"
	"net/http"
	"sync"
)

//todo
// remove file extention from name searchText
// take length into acount. ie searching "jacks" should find "jacks books" before "bills"

const limit int = 15
const maxDist int = 25

func levelshteinDist(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// Make sure we use the shorter string for the inner loop
	if la < lb {
		// swap to ensure lb <= la
		a, b = b, a
		la, lb = lb, la
	}
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	// initialize row 0
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		ai := a[i-1]
		for j := 1; j <= lb; j++ {
			cost := 0
			if ai != b[j-1] {
				cost = 1
			}
			// substitution, insertion, deletion
			sub := prev[j-1] + cost
			ins := curr[j-1] + 1
			del := prev[j] + 1
			// take min
			if ins < sub {
				sub = ins
			}
			if del < sub {
				sub = del
			}
			curr[j] = sub
		}
		// swap rows for next iteration
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

	for currentRankedFile := range resultChan {

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

	return res
}

func exploreFolder(f *Folder, text string, c chan<- rankedFile, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, folder := range f.Subfolders {
		wg.Add(1)
		go exploreFolder(folder, text, c, wg)
	}
	for _, file := range f.Files {
		dist := levelshteinDist(file.Name, text)
		if dist <= maxDist {

			c <- rankedFile{file: *file, distance: dist}
		}
	}
}
