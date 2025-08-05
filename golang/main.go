package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/COS301-SE-2025/Smart-File-Manager/golang/filesystem"
)

// RAKE stopwords
var stopwords = map[string]struct{}{
	"and": {}, "the": {}, "is": {}, "in": {}, "at": {}, "of": {}, "a": {}, "an": {},
	"to": {}, "for": {}, "on": {}, "with": {}, "as": {}, "by": {}, "that": {}, "this": {}, "if": {},
}

// Keyword holds a word and its RAKE‐derived score
type Keyword struct {
	Word  string
	Score float64
}

// ExtractKeywordsRAKE reads a single file at filePath (skipping binaries),
// runs RAKE on its text, and returns the topN words by score.
func ExtractKeywordsRAKE(filePath string, topN int) ([]Keyword, error) {
	// 1) open & sniff
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	f.Seek(0, io.SeekStart)
	if !strings.HasPrefix(http.DetectContentType(buf[:n]), "text/") {
		return nil, fmt.Errorf("not a text file: %s", filePath)
	}

	// 2) read all text
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	text := string(data)

	// 3) tokenize into words
	wordRe := regexp.MustCompile(`[A-Za-z0-9]+`)
	words := wordRe.FindAllString(strings.ToLower(text), -1)

	// 4) build candidate phrases
	var phrases [][]string
	var cur []string
	for _, w := range words {
		if _, isStop := stopwords[w]; isStop {
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

	// 5) compute word freq & degree
	freq := map[string]float64{}
	degree := map[string]float64{}
	for _, phrase := range phrases {
		L := float64(len(phrase))
		for _, w := range phrase {
			freq[w]++
			degree[w] += L
		}
	}

	// 6) score = degree / freq
	score := map[string]float64{}
	for w, f := range freq {
		score[w] = degree[w] / f
	}

	// 7) collect & sort topN
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
	return kws, nil
}

func main() {

	keywords, err := ExtractKeywordsRAKE("main.go", 20)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	for _, kw := range keywords {
		fmt.Printf("%-15s %.2f\n", kw.Word, kw.Score)
	}
	//testing directory creation
	// folder := mockFolderStructure()
	// filesystem.CreateDirectoryStructure(folder)
	filesystem.HandleRequests()

	// print current Composites
	Composites := filesystem.GetComposites()
	for _, item := range Composites {
		item.Display(0)
	}
	// print current composites
	// composites := filesystem.GetComposites()
	// for _, item := range composites {
	// 	item.Display(0)
	// }
}
func mockFolderStructure() *filesystem.Folder {
	return &filesystem.Folder{
		NewPath: "test_root",
		Subfolders: []*filesystem.Folder{
			{
				NewPath: "test_root/sub1",
				Subfolders: []*filesystem.Folder{
					{
						NewPath: "test_root/sub1/sub1_1",
					},
				},
			},
			{
				NewPath: "test_root/sub2",
			},
		},
	}
}
