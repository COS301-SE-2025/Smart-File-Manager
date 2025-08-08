package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func keywordSearchHadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")

	for _, c := range Composites {
		if c.Name == name {
			// build the nested []FileNode
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
			for _, file := range c.Files {

				keywords, err := ExtractKeywordsRAKE(ConvertToWSLPath(file.Path), 20, 50*1024*1024)
				// keywords, err := ExtractKeywordsRAKE("document.txt", 20)

				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					os.Exit(1)
				}
				for _, kw := range keywords {
					fmt.Printf("%-15s %.2f\n", kw.Word, kw.Score)
				}
			}

			goElapsed := time.Since(goStart)

			fmt.Printf("grpc Code block executed in %s\n", grpcElapsed)
			fmt.Printf("go Code block executed in %s\n", goElapsed)

			return
		}
	}

	http.Error(w, "No smart manager with that name", http.StatusBadRequest)
}

var stopwords = map[string]struct{}{
	"and": {}, "the": {}, "is": {}, "in": {}, "at": {}, "of": {}, "a": {}, "an": {},
	"to": {}, "for": {}, "on": {}, "with": {}, "as": {}, "by": {}, "that": {}, "this": {}, "if": {},
}

// Keyword holds a word and its RAKE score
type Keyword struct {
	Word  string
	Score float64
}

// ExtractKeywordsFromText runs RAKE on the given text and returns topN words.
func ExtractKeywordsFromText(text string, topN int) []Keyword {
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

	freq := map[string]float64{}
	degree := map[string]float64{}
	for _, p := range phrases {
		L := float64(len(p))
		for _, w := range p {
			freq[w]++
			degree[w] += L
		}
	}

	score := map[string]float64{}
	for w, f := range freq {
		score[w] = degree[w] / f
	}

	var kws []Keyword
	for w, sc := range score {
		kws = append(kws, Keyword{Word: w, Score: sc})
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
func ExtractKeywordsRAKE(filePath string, topN int, maxSize int64) ([]Keyword, error) {
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

	data, err := ioutil.ReadAll(io.LimitReader(f, maxSize))
	if err != nil {
		return nil, nil
	}

	return ExtractKeywordsFromText(string(data), topN), nil
}
