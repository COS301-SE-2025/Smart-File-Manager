package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
)

//flow idea:
// app starts
// keywords are loaded from json(and locks and tags) if they exist
// run go version of keyword extraction and python at the same time
// overwrite keywords for all files when go returns
// overwrite keywords for all files when python returns
// save python keywords along with tags and locks

var wg sync.WaitGroup

func keywordSearchHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")

	for _, c := range Composites {
		if c.Name == name {

			grpcStart := time.Now()
			err := grpcFunc(c, "KEYWORDS")
			if err != nil {
				log.Fatalf("grpcFunc failed: %v", err)
				http.Error(w, "internal server error, GRPC CALLED WRONG", http.StatusInternalServerError)
			}
			children := createDirectoryJSONStructure(c)

			root := DirectoryTreeJson{
				Name:     c.Name,
				IsFolder: true,
				Children: children,
			}

			if err := json.NewEncoder(w).Encode(root); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			grpcElapsed := time.Since(grpcStart)

			goStart := time.Now()
			extractKeywords(c)

			goElapsed := time.Since(goStart)

			PrettyPrintFolder(c, "")

			fmt.Printf("grpc Code block executed in %s\n", grpcElapsed)
			fmt.Printf("go Code block executed in %s\n", goElapsed)

			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

func extractKeywords(c *Folder) {
	wg.Add(1)
	go getKeywords(c)
	wg.Wait()
}

func getKeywords(c *Folder) {
	defer wg.Done()

	for _, file := range c.Files {

		keywords, err := ExtractKeywordsRAKE(ConvertToWSLPath(file.Path), 20, 50*1024*1024)

		if err != nil {
			fmt.Printf("err")
		}
		mu.Lock()
		file.Keywords = keywords
		mu.Unlock()
	}

	for _, folder := range c.Subfolders {
		wg.Add(1)
		go getKeywords(folder)
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
