package filesystem

import (
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

func keywordSearchHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")

	for _, c := range Composites {
		if c.Name == name {
			// pythonExtractKeywords(c)
			goExtractKeywords(c)
			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

func pythonExtractKeywords(c *Folder) {
	fmt.Println("started")
	grpcStart := time.Now()
	err := grpcFunc(c, "KEYWORDS")
	if err != nil {
		log.Fatalf("grpcFunc failed: %v", err)
	}

	saveCompositeDetails(c)

	fmt.Println("finished python")

	grpcElapsed := time.Since(grpcStart)
	fmt.Printf("grpc Code block executed in %s\n", grpcElapsed)

}

func goExtractKeywords(c *Folder) {
	var wg sync.WaitGroup

	// IMPORTANT: do NOT start this as a goroutine.
	// We need all wg.Adds to happen before Wait begins.
	getKeywords(c, &wg)

	wg.Wait()
	saveCompositeDetails(c)
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
