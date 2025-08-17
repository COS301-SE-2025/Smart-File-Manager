package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
)

const limitKeywordSearch int = 20
const maxDistKeywordSearch int = 3

//flow idea:
// app starts
// keywords are loaded from json(and locks and tags) if they exist
// run go version of keyword extraction and python at the same time
// overwrite keywords for all files when go returns
// overwrite keywords for all files when python returns
// save python keywords along with tags and locks

// route to check if keywords have been populated
func IsKeywordSearchReadyHander(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("compositeName")

	for _, c := range Composites {
		if c.Name == name {
			hasKeywords := c.HasKeywords
			if err := json.NewEncoder(w).Encode(hasKeywords); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

func KeywordSearchHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify your frontend origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	name := r.URL.Query().Get("compositeName")
	searchText := r.URL.Query().Get("searchText")

	for _, c := range Composites {
		if c.Name == name {

			var terms []string = strings.Split(searchText, " ")

			sr := getMatchesByKeywords(terms, c)

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

func getMatchesByKeywords(searchTerms []string, composite *Folder) *safeResults {

	res := &safeResults{
		Name: composite.Name,
	}

	resultChan := make(chan rankedFile)

	var wg sync.WaitGroup

	wg.Add(1)
	go exploreFolderForKeywords(composite, searchTerms, resultChan, &wg)

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

	return res
}

func exploreFolderForKeywords(f *Folder, searchTerms []string, c chan<- rankedFile, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, folder := range f.Subfolders {
		wg.Add(1)
		go exploreFolderForKeywords(folder, searchTerms, c, wg)
	}

	for _, file := range f.Files {
		finalDist := 100
		for _, searchTerm := range searchTerms {
			for _, keyword := range file.Keywords {

				dist := LevenshteinDistForKeywords(searchTerm, keyword.Keyword)
				if dist <= maxDistKeywordSearch {
					if dist < finalDist {
						finalDist = dist
						fmt.Println("updating dist for FILE: " + file.Name + " with dist: " + strconv.Itoa(finalDist))
						fmt.Println("distance between: " + searchTerm + " AND " + keyword.Keyword + " = " + strconv.Itoa(dist))
					}

				}
			}
		}
		if finalDist <= maxDistKeywordSearch {
			c <- rankedFile{file: *file, distance: finalDist}
		}
	}
}

func LevenshteinDistForKeywords(searchText string, fileKeyword string) int {
	if len(searchText) == 0 {
		return len(fileKeyword)
	}
	if len(fileKeyword) == 0 {
		return len(searchText)
	}

	searchText = strings.ToLower(searchText)
	fileKeyword = strings.ToLower(fileKeyword)

	//  exact matches should be 0
	if fileKeyword == searchText {
		return 0
	}
	//this is the name
	//his ist the
	//  BOOST exact substrings
	var boost float32 = 1 //lower boost is better as it makes the distance smaller
	if strings.Contains(fileKeyword, searchText) {
		boost = 0.1
	}

	// now fall back on full Levenshtein
	lenSearchText, lenFileName := len(searchText), len(fileKeyword)

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
			if ai != fileKeyword[j-1] {
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

func pythonExtractKeywords(c *Folder) {
	err := grpcFunc(c, "KEYWORDS")
	if err != nil {
		log.Fatalf("grpcFunc failed: %v", err)
	}

	saveCompositeDetails(c)
	c.HasKeywords = true

}

func GoExtractKeywords(c *Folder) {
	var wg sync.WaitGroup

	// IMPORTANT: do NOT start this as a goroutine.
	// We need all wg.Adds to happen before Wait begins.
	getKeywords(c, &wg)

	wg.Wait()
	saveCompositeDetails(c)
	c.HasKeywords = true
}

var keywordSem = make(chan struct{}, 32) // tune limit

func getKeywords(c *Folder, wg *sync.WaitGroup) {
	for _, file := range c.Files {
		f := file // capture loop var (or pass as param)
		wg.Add(1)
		go func(f *File) {
			defer wg.Done()
			keywordSem <- struct{}{}
			defer func() { <-keywordSem }()

			kw, err := ExtractKeywordsRAKE(
				ConvertToWSLPath(f.Path),
				20,
				50*1024*1024,
			)
			if err != nil {
				fmt.Printf("extract keywords error for %s: %v\n", f.Path, err)
				return
			}
			f.Keywords = kw
		}(f)
	}

	// recurse synchronously so all Adds happen before Wait
	for _, folder := range c.Subfolders {
		getKeywords(folder, wg)
	}
}

var stopwords = map[string]struct{}{
	"and": {}, "the": {}, "is": {}, "in": {}, "at": {}, "of": {}, "a": {}, "an": {},
	"to": {}, "for": {}, "on": {}, "with": {}, "as": {}, "by": {}, "that": {}, "this": {}, "if": {},
}

// Keyword holds a word and its RAKE score

// ExtractKeywordsFromText runs RAKE on the given text and returns topN words.
func ExtractKeywordsFromText(text string, topN int) []*pb.Keyword {
	re := regexp.MustCompile(`[A-Za-z0-9]+`)
	words := re.FindAllString(strings.ToLower(text), -1)

	var phrases [][]string
	var cur []string
	for _, w := range words {
		if _, stop := stopwords[w]; stop {
			if len(cur) > 0 {
				phrases = append(phrases, cur)
				cur = nil
			}
		} else {
			cur = append(cur, w)
		}
	}
	if len(cur) > 0 {
		phrases = append(phrases, cur)
	}

	freq := map[string]float32{}
	degree := map[string]float32{}
	for _, p := range phrases {
		L := float32(len(p))
		for _, w := range p {
			freq[w]++
			degree[w] += L
		}
	}

	score := map[string]float32{}
	for w, f := range freq {
		score[w] = degree[w] / f
	}

	var kws []*pb.Keyword
	for w, sc := range score {
		kws = append(kws, &pb.Keyword{Keyword: w, Score: sc})
	}
	sort.Slice(kws, func(i, j int) bool {
		return kws[i].Score > kws[j].Score
	})
	if len(kws) > topN {
		kws = kws[:topN]
	}
	return kws
}

// ExtractKeywordsRAKE reads filePath (limit maxSize), but only if it's a text/* file.
// Otherwise it returns an empty slice.
func ExtractKeywordsRAKE(filePath string, topN int, maxSize int64) ([]*pb.Keyword, error) {
	fi, err := os.Stat(filePath)
	if err != nil || fi.IsDir() || fi.Size() > maxSize {
		return nil, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	head := make([]byte, 512)
	n, _ := f.Read(head)
	f.Seek(0, io.SeekStart)

	ctype := http.DetectContentType(head[:n])
	if !strings.HasPrefix(ctype, "text/") {
		// not a plain-text file → no keywords
		return nil, nil
	}

	data, err := io.ReadAll(io.LimitReader(f, maxSize))
	if err != nil {
		return nil, nil
	}

	return ExtractKeywordsFromText(string(data), topN), nil
}
