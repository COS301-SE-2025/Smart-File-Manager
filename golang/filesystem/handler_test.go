package filesystem

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
)

// stub out grpcFunc so handler tests don’t talk to real server
func init() {
	grpcFunc = func(c *Folder, requestType string) error {
		return nil
	}
}

// --- helpers to build sample data ---

func makeSampleProtoDir() *pb.Directory {
	return &pb.Directory{
		Name: "pdir",
		Path: "/pdir",
		Files: []*pb.File{
			{Name: "f1", OriginalPath: "/pdir/f1", Tags: []*pb.Tag{{Name: "t1"}}},
		},
		Directories: []*pb.Directory{
			{Name: "sub", Path: "/pdir/sub"},
		},
	}
}

func makeSampleFolder() *Folder {
	return &Folder{
		Name: "fdir",
		Path: "/fdir",
		Files: []*File{
			{Name: "f1", Path: "/fdir/f1", Tags: []string{"t1"}, Metadata: []*MetadataEntry{{Key: "size_bytes", Value: "99"}}},
		},
		Subfolders: []*Folder{{Name: "sub", Path: "/fdir/sub"}},
	}
}

// --- unit tests ---

func TestTagsConversion(t *testing.T) {
	in := []*pb.Tag{{Name: "a"}, {Name: "b"}}
	got := tagsToStrings(in)
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tagsToStrings = %v, want %v", got, want)
	}

	sin := []string{"x", "y"}
	stags := stringsToTags(sin)
	if len(stags) != 2 || stags[0].Name != "x" || stags[1].Name != "y" {
		t.Errorf("stringsToTags = %+v, want pb.Tags x,y", stags)
	}
}

func TestMetadataConverterAndExtract(t *testing.T) {
	pbs := []*pb.MetadataEntry{
		{Key: "size_bytes", Value: "10"},
		{Key: "created", Value: "2025-01-02T15:04:05Z"},
	}
	me := metadataConverter(pbs)
	if len(me) != 2 || me[0].Key != "size_bytes" || me[1].Value != "2025-01-02T15:04:05Z" {
		t.Fatalf("metadataConverter = %+v", me)
	}

	md := extractMetadata(me)
	if md.Size != "10" {
		t.Errorf("extractMetadata.Size = %q, want %q", md.Size, "10")
	}
	if md.DateCreated != "2025-01-02T15:04:05Z" {
		t.Errorf("extractMetadata.DateCreated = %q", md.DateCreated)
	}
}

func TestConvertFolderProtoRoundTrip(t *testing.T) {
	orig := makeSampleFolder()
	proto := convertFolderToProto(*orig)
	round := convertProtoToFolder(proto)

	if round.Name != orig.Name || round.Path != orig.Path {
		t.Errorf("roundtrip Name/Path = %v/%v, want %v/%v",
			round.Name, round.Path, orig.Name, orig.Path)
	}
	// tags preserved?
	if len(round.Files) != 1 || len(round.Files[0].Tags) != 1 || round.Files[0].Tags[0] != "t1" {
		t.Errorf("roundtrip file tags = %v", round.Files[0].Tags)
	}
}

func TestCreateDirectoryJSONStructure(t *testing.T) {
	f := makeSampleFolder()
	nodes := createDirectoryJSONStructure(f)

	// expect two nodes: the file, then the subfolder
	if len(nodes) != 2 {
		t.Fatalf("got %d nodes, want 2", len(nodes))
	}
	if nodes[0].Name != "f1" || nodes[1].Name != "sub" || !nodes[1].IsFolder {
		t.Errorf("nodes = %+v", nodes)
	}
	// sub should have no children
	if len(nodes[1].Children) != 0 {
		t.Errorf("sub.Children = %v, want empty", nodes[1].Children)
	}
}

// --- integration‐style test for the HTTP handler ---

func TestLoadTreeDataHandler(t *testing.T) {
	// prepare
	composites = []*Folder{makeSampleFolder()}
	composites[0].Name = "foo"

	req := httptest.NewRequest("GET", "/loadTreeData?name=foo", nil)
	w := httptest.NewRecorder()

	loadTreeDataHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	var out DirectoryTreeJson
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.Name != "foo" {
		t.Errorf("out.Name = %q, want foo", out.Name)
	}
	if len(out.Children) != 2 {
		t.Errorf("out.Children len = %d, want 2", len(out.Children))
	}
}

// --- a tiny test for the printer: it should not panic on nil and print something predictable ---

func TestPrintFolderDetailsNil(t *testing.T) {
	// should just return, not panic
	printFolderDetails(nil, 3)
}

func TestPrintFolderDetailsOutput(t *testing.T) {
	f := &Folder{
		Name:         "X",
		Path:         "/p",
		newPath:      "/np",
		CreationDate: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
		Locked:       true,
		Tags:         []string{"a"},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printFolderDetails(f, 0)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "Folder: X") {
		t.Errorf("output = %q; want it to contain %q", out, "Folder: X")
	}
}
