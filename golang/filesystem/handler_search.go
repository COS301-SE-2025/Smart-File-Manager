package filesystem

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const limit int = 25
const maxDist int = 14
const similarityThreshold = 0.4

func LevenshteinDist(searchText string, fileName string) int {
	if len(searchText) == 0 {
		return len(fileName)
	}
	if len(fileName) == 0 {
		return len(searchText)
	}

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

	//  exact matches should be 0
	if fileName == searchText {
		return 0
	}
	//this is the name
	//his ist the
	//  BOOST exact substrings
	var boost float32 = 1 //lower boost is better as it makes the distance smaller
	if strings.Contains(fileName, searchText) {
		boost = 0.1
	}

	// now fall back on full Levenshtein
	lenSearchText, lenFileName := len(searchText), len(fileName)

	prev := make([]int, lenFileName+1)
	curr := make([]int, lenFileName+1)
	for j := 0; j <= lenFileName; j++ {
		prev[j] = j
	}
	for i := 1; i <= lenSearchText; i++ {
		curr[0] = i
		ai := searchText[i-1]
		for j := 1; j <= lenFileName; j++ {
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

	return int(math.Round(float64(prev[lenFileName]) * float64(boost)))

}

type safeResults struct {
	Name        string
	rankedFiles []rankedFile
}

type rankedFile struct {
	file     File
	distance int
}

// gpt given
func PrettyPrintFolder(f *Folder, indent string) {
	fmt.Printf("%s📁 %s (path = %v)\n", indent, f.Name, f.Path)
	// print files
	for _, file := range f.Files {
		fmt.Printf("%s  📄 %s\n", indent, file.Name)
		for _, tag := range file.Tags {
			fmt.Printf("%s  TAG: %s\n", indent, tag)
		}
		// fmt.Printf("%s  locked: %s\n", indent, file.Locked)

		fmt.Println("----")
	}
	// recurse into subfolders
	for _, sub := range f.Subfolders {
		PrettyPrintFolder(sub, indent+"  ")
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	compositeName := r.URL.Query().Get("compositeName")
	searchText := r.URL.Query().Get("searchText")

	for _, comp := range Composites {
		if comp.Name == compositeName {

			sr := getMatches(searchText, comp)

			cores := DirectoryTreeJson{
				Name:     sr.Name,
				IsFolder: true,
				Children: make([]FileNode, len(sr.rankedFiles)),
			}

			for i, rf := range sr.rankedFiles {
				file := rf.file
				fileNodeToAdd := FileNode{
					Name:     file.Name,
					Path:     file.Path,
					IsFolder: false,
					Tags:     file.Tags,
					Metadata: ConvertMetadataEntries(file.Metadata),
				}
				cores.Children[i] = fileNodeToAdd
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cores)
			return
		}
	}
	http.Error(w, "No smart manager with that name", http.StatusBadRequest)

}

func ConvertMetadataEntries(entries []*MetadataEntry) *Metadata {
	md := &Metadata{}

	for _, entry := range entries {
		switch strings.ToLower(entry.Key) {

		case "size":
			md.Size = entry.Value
		case "datecreated":
			md.DateCreated = entry.Value
		case "lastmodified":
			md.LastModified = entry.Value
		case "mimetype":
			md.MimeType = entry.Value

		}
	}

	return md
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

				if len(res.rankedFiles) < limit && (currentRankedFile.distance < maxDist) {
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
			if !inserted && (currentRankedFile.distance < maxDist) {
				res.rankedFiles = append(res.rankedFiles, currentRankedFile)
			}
		}
	}

	//checks to remove dups (yes if concurrency was perfect there wouldnt be dups)
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

	sort.SliceStable(res.rankedFiles, func(i, j int) bool {
		if res.rankedFiles[i].distance != res.rankedFiles[j].distance {
			return res.rankedFiles[i].distance < res.rankedFiles[j].distance
		}
		return strings.ToLower(res.rankedFiles[i].file.Name) < strings.ToLower(res.rankedFiles[j].file.Name)
	})

	return res
}

func exploreFolder(f *Folder, text string, c chan<- rankedFile, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, folder := range f.Subfolders {
		wg.Add(1)
		go exploreFolder(folder, text, c, wg)
	}

	for _, file := range f.Files {
		dist := LevenshteinDist(text, file.Name)
		maxLen := len(text)
		if len(file.Name) > maxLen {
			maxLen = len(file.Name)
		}
		// avoid division by zero
		if maxLen == 0 {
			continue
		}
		similarity := 1.0 - float64(dist)/float64(maxLen)

		if dist <= maxDist && similarity >= similarityThreshold {
			c <- rankedFile{file: *file, distance: dist}
		}
	}
}
