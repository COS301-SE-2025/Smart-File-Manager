package filesystem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	pb "github.com/COS301-SE-2025/Smart-File-Manager/golang/client/protos"
)

func TestMoveDirectoryHandler_Success(t *testing.T) {
	setupTest(t)
	defer cleanupTest(t)

	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)

	testFile := filepath.Join(sourceDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	testComposite := &Folder{
		Name: "testManager",
		Path: sourceDir,
		Files: []*File{
			{
				Name:    "test.txt",
				Path:    testFile,
				NewPath: "testManager/test.txt",
			},
		},
	}

	originalComposites := Composites
	originalObjectMap := ObjectMap
	Composites = []*Folder{testComposite}
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[sourceDir] = make(map[string]object)
	defer func() {
		Composites = originalComposites
		ObjectMap = originalObjectMap
	}()

	req := httptest.NewRequest("GET", "/move?name=testManager", nil)
	w := httptest.NewRecorder()

	moveDirectoryHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	response := strings.TrimSpace(w.Body.String())
	if response != "true" {
		t.Errorf("Expected 'true', got %s", response)
	}

	// Check composite count
	if len(Composites) != 1 {
		t.Errorf("Expected 1 composite after move, got %d", len(Composites))
	}

	// Check new directory exists
	expectedDir := filepath.Join(tempDir, "testManager")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected new directory %s to exist after move", expectedDir)
	}

	// Check old directory is removed
	if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
		t.Errorf("Expected old directory %s to be removed", sourceDir)
	}

	// Check composite path is updated
	if Composites[0].Path != expectedDir {
		t.Errorf("Expected composite path %s, got %s", expectedDir, Composites[0].Path)
	}
}

func TestMoveDirectoryHandler_NotFound(t *testing.T) {
	setupTest(t)
	defer cleanupTest(t)

	originalComposites := Composites
	Composites = []*Folder{}
	defer func() { Composites = originalComposites }()

	req := httptest.NewRequest("GET", "/move?name=nonexistent", nil)
	w := httptest.NewRecorder()

	moveDirectoryHandler(w, req)

	response := strings.TrimSpace(w.Body.String())
	if response != "false" {
		t.Errorf("Expected 'false', got %s", response)
	}
}

func TestMoveContent(t *testing.T) {
	setupTest(t)
	defer cleanupTest(t)

	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)

	testFile := filepath.Join(sourceDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	item := &Folder{
		Name: "testManager",
		Path: sourceDir,
		Files: []*File{
			{
				Name:    "test.txt",
				Path:    testFile,
				NewPath: "organized/test.txt",
			},
		},
	}

	CreateDirectoryStructure(item)
	moveContent(item)

	expectedPath := filepath.Join(tempDir, "testManager")
	if item.Path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, item.Path)
	}

	if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
		t.Error("Expected original directory to be removed")
	}
}

func TestMoveContentRecursive(t *testing.T) {
	tempDir := t.TempDir()
	root = tempDir

	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)

	testFile := filepath.Join(sourceDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	subSourceDir := filepath.Join(sourceDir, "subdir")
	os.MkdirAll(subSourceDir, 0755)

	subFile := filepath.Join(subSourceDir, "sub.txt")
	os.WriteFile(subFile, []byte("subcontent"), 0644)

	item := &Folder{
		Name: "parent",
		Files: []*File{
			{
				Name:    "test.txt",
				Path:    testFile,
				NewPath: "moved/test.txt",
			},
		},
		Subfolders: []*Folder{
			{
				Name:    "subdir",
				NewPath: "moved/subdir",
				Files: []*File{
					{
						Name:    "sub.txt",
						Path:    subFile,
						NewPath: "moved/subdir/sub.txt",
					},
				},
			},
		},
	}

	os.MkdirAll(filepath.Join(tempDir, "moved"), 0755)
	os.MkdirAll(filepath.Join(tempDir, "moved", "subdir"), 0755)

	moveContentRecursive(item)

	movedFile := filepath.Join(tempDir, "moved", "test.txt")
	if _, err := os.Stat(movedFile); os.IsNotExist(err) {
		t.Error("Expected moved file to exist")
	}

	movedSubFile := filepath.Join(tempDir, "moved", "subdir", "sub.txt")
	if _, err := os.Stat(movedSubFile); os.IsNotExist(err) {
		t.Error("Expected moved subfolder file to exist")
	}

	expectedSubfolderPath := filepath.Join(tempDir, "moved/subdir")
	if item.Subfolders[0].Path != expectedSubfolderPath {
		t.Errorf("Expected subfolder path %s, got %s", expectedSubfolderPath, item.Subfolders[0].Path)
	}
}

func TestMoveContentRecursive_NilFolder(t *testing.T) {
	moveContentRecursive(nil)
}

func TestGenerateUniqueFilePath(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentPath := filepath.Join(tempDir, "new.txt")
	result := generateUniqueFilePath(nonExistentPath)
	if result != nonExistentPath {
		t.Errorf("Expected %s, got %s", nonExistentPath, result)
	}

	existingFile := filepath.Join(tempDir, "existing.txt")
	os.WriteFile(existingFile, []byte("content"), 0644)

	uniquePath := generateUniqueFilePath(existingFile)
	expected := filepath.Join(tempDir, "existing_(1).txt")
	if uniquePath != expected {
		t.Errorf("Expected %s, got %s", expected, uniquePath)
	}

	os.WriteFile(expected, []byte("content"), 0644)
	uniquePath2 := generateUniqueFilePath(existingFile)
	expected2 := filepath.Join(tempDir, "existing_(2).txt")
	if uniquePath2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, uniquePath2)
	}
}

func TestGenerateUniqueFilePath_NoExtension(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "noext")
	os.WriteFile(existingFile, []byte("content"), 0644)

	uniquePath := generateUniqueFilePath(existingFile)
	expected := filepath.Join(tempDir, "noext_(1)")
	if uniquePath != expected {
		t.Errorf("Expected %s, got %s", expected, uniquePath)
	}
}

func TestCreateDirectoryStructure(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test_manager_path")
	os.MkdirAll(testPath, 0755)

	folder := &Folder{
		Name:    "manager1",
		Path:    testPath,
		NewPath: "test_root",
		Subfolders: []*Folder{
			{
				Name:    "sub1",
				NewPath: "test_root/sub1",
				Subfolders: []*Folder{
					{
						Name:    "sub1_1",
						NewPath: "test_root/sub1/sub1_1",
					},
				},
			},
			{
				Name:    "sub2",
				NewPath: "test_root/sub2",
			},
		},
	}

	CreateDirectoryStructure(folder)

	managerRoot := filepath.Join(testPath, "manager1")
	var got []string
	err := filepath.Walk(managerRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			rel, err := filepath.Rel(managerRoot, path)
			if err != nil {
				return err
			}
			got = append(got, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking manager directory: %v", err)
	}

	want := []string{
		".",
		"test_root",
		"test_root/sub1",
		"test_root/sub1/sub1_1",
		"test_root/sub2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("directory tree mismatch:\n got:  %#v\n want: %#v", got, want)
	}
}

func TestCreateDirectoryStructureRecursive_NilFolder(t *testing.T) {
	CreateDirectoryStructureRecursive(nil)
}

func TestCreateDirectoryStructureRecursive_EmptySubfolders(t *testing.T) {
	tempDir := t.TempDir()
	root = tempDir

	folder := &Folder{
		Name:    "empty",
		NewPath: "empty",
	}

	CreateDirectoryStructureRecursive(folder)

	expectedPath := filepath.Join(tempDir, "empty")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}

func TestUpdateStoredPathsFromComposite_NoExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, "storage")
	os.MkdirAll(storageDir, 0755)

	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	comp := &Folder{
		Name: "test",
		Path: "/test/path",
		Files: []*File{
			{
				Name: "file1.txt",
				Path: "/test/path/file1.txt",
				Tags: []string{"tag1"},
			},
		},
	}

	err := UpdateStoredPathsFromComposite(comp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	filePath := filepath.Join("storage", "test.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected storage file to be created")
	}
}

func TestUpdateStoredPathsFromComposite_WithExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, "storage")
	os.MkdirAll(storageDir, 0755)

	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	oldStructure := DirectoryTreeJson{
		Name:     "test",
		IsFolder: true,
		RootPath: "/old/path",
		Children: []FileNode{
			{
				Name:     "file1.txt",
				Path:     "/old/path/file1.txt",
				IsFolder: false,
				Tags:     []string{"oldtag"},
				Keywords: []*pb.Keyword{{Keyword: "oldkeyword"}},
				Locked:   true,
			},
		},
	}

	filePath := filepath.Join("storage", "test.json")
	data, _ := json.MarshalIndent(oldStructure, "", "  ")
	os.WriteFile(filePath, data, 0644)

	comp := &Folder{
		Name: "test",
		Path: "/new/path",
		Files: []*File{
			{
				Name: "file1.txt",
				Path: "/new/path/file1.txt",
			},
		},
	}

	err := UpdateStoredPathsFromComposite(comp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	newData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	var newStructure DirectoryTreeJson
	json.Unmarshal(newData, &newStructure)

	if newStructure.RootPath != "/new/path" {
		t.Errorf("Expected root path /new/path, got %s", newStructure.RootPath)
	}

	if len(newStructure.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(newStructure.Children))
	}

	child := newStructure.Children[0]
	if child.Path != "/new/path/file1.txt" {
		t.Errorf("Expected path /new/path/file1.txt, got %s", child.Path)
	}

	if len(child.Tags) == 0 || child.Tags[0] != "oldtag" {
		t.Errorf("Expected preserved tag 'oldtag', got %v", child.Tags)
	}

	if !child.Locked {
		t.Error("Expected locked status to be preserved")
	}
}

func TestBuildNodesWithPreservedMetadata(t *testing.T) {
	oldStructure := &DirectoryTreeJson{
		Children: []FileNode{
			{
				Name:     "file1.txt",
				Path:     "/old/file1.txt",
				IsFolder: false,
				Tags:     []string{"tag1"},
				Keywords: []*pb.Keyword{{Keyword: "keyword1"}},
				Locked:   true,
			},
			{
				Name:     "folder1",
				Path:     "/old/folder1",
				IsFolder: true,
				Tags:     []string{"foldertag"},
				Locked:   true,
				Children: []FileNode{
					{
						Name:     "nested.txt",
						Path:     "/old/folder1/nested.txt",
						IsFolder: false,
						Tags:     []string{"nested"},
					},
				},
			},
		},
	}

	folder := &Folder{
		Files: []*File{
			{
				Name: "file1.txt",
				Path: "/new/file1.txt",
			},
			{
				Name: "newfile.txt",
				Path: "/new/newfile.txt",
				Tags: []string{"newtag"},
			},
		},
		Subfolders: []*Folder{
			{
				Name: "folder1",
				Path: "/new/folder1",
				Files: []*File{
					{
						Name: "nested.txt",
						Path: "/new/folder1/nested.txt",
					},
				},
			},
		},
	}

	nodes := buildNodesWithPreservedMetadata(folder, oldStructure)

	if len(nodes) != 3 {
		t.Fatalf("Expected 3 nodes, got %d", len(nodes))
	}

	file1Found := false
	folderFound := false
	for _, node := range nodes {
		if node.Name == "file1.txt" {
			file1Found = true
			if len(node.Tags) == 0 || node.Tags[0] != "tag1" {
				t.Error("Expected preserved tag for file1.txt")
			}
			if !node.Locked {
				t.Error("Expected preserved lock status for file1.txt")
			}
			if len(node.Keywords) == 0 || node.Keywords[0].Keyword != "keyword1" {
				t.Error("Expected preserved keyword for file1.txt")
			}
		}
		if node.Name == "folder1" && node.IsFolder {
			folderFound = true
			if len(node.Tags) == 0 || node.Tags[0] != "foldertag" {
				t.Error("Expected preserved tag for folder1")
			}
			if !node.Locked {
				t.Error("Expected preserved lock status for folder1")
			}
		}
	}

	if !file1Found {
		t.Error("file1.txt node not found")
	}
	if !folderFound {
		t.Error("folder1 node not found")
	}
}

func TestBuildPathMap(t *testing.T) {
	nodes := []FileNode{
		{
			Name:     "file1.txt",
			Path:     "/path/file1.txt",
			IsFolder: false,
		},
		{
			Name:     "folder1",
			Path:     "/path/folder1",
			IsFolder: true,
			Children: []FileNode{
				{
					Name:     "nested.txt",
					Path:     "/path/folder1/nested.txt",
					IsFolder: false,
				},
			},
		},
	}

	pathMap := make(map[string]FileNode)
	buildPathMap(nodes, pathMap)

	expectedEntries := []string{
		"/path/file1.txt",
		"file1.txt",
		"/path/folder1",
		"folder1",
		"/path/folder1/nested.txt",
		"nested.txt",
	}

	for _, key := range expectedEntries {
		if _, exists := pathMap[key]; !exists {
			t.Errorf("Expected key %s to exist in pathMap", key)
		}
	}

	if len(pathMap) != len(expectedEntries) {
		t.Errorf("Expected %d entries in pathMap, got %d", len(expectedEntries), len(pathMap))
	}
}

func TestFindNodeByName(t *testing.T) {
	pathMap := map[string]FileNode{
		"file1.txt": {
			Name:     "file1.txt",
			IsFolder: false,
			Tags:     []string{"tag1"},
		},
		"folder1": {
			Name:     "folder1",
			IsFolder: true,
			Tags:     []string{"foldertag"},
		},
	}

	node, found := findNodeByName(pathMap, "file1.txt", false)
	if !found {
		t.Error("Expected to find file1.txt")
	}
	if node.Name != "file1.txt" {
		t.Errorf("Expected name file1.txt, got %s", node.Name)
	}

	_, found = findNodeByName(pathMap, "file1.txt", true)
	if found {
		t.Error("Should not find file1.txt as folder")
	}

	folderNode, found := findNodeByName(pathMap, "folder1", true)
	if !found {
		t.Error("Expected to find folder1")
	}
	if folderNode.Name != "folder1" {
		t.Errorf("Expected name folder1, got %s", folderNode.Name)
	}

	_, found = findNodeByName(pathMap, "nonexistent", false)
	if found {
		t.Error("Should not find nonexistent file")
	}
}

func TestGetPath(t *testing.T) {
	setupTest(t)
	defer cleanupTest(t)

	path := getPath()
	if !strings.Contains(path, "Smart-File-Manager") {
		t.Errorf("Expected path to contain Smart-File-Manager, got %s", path)
	}
}

func TestCleanManagerPrefix(t *testing.T) {
	tests := []struct {
		path        string
		managerName string
		expected    string
	}{
		{
			path:        "/home/user/manager1/manager1/file.txt",
			managerName: "manager1",
			expected:    "/home/user/manager1/file.txt",
		},
		{
			path:        "/home/manager1/folder/manager1/manager1/file.txt",
			managerName: "manager1",
			expected:    "/home/manager1/folder/file.txt",
		},
		{
			path:        "/simple/path/file.txt",
			managerName: "manager1",
			expected:    "/simple/path/file.txt",
		},
	}

	for _, test := range tests {
		result := cleanManagerPrefix(test.path, test.managerName)
		if result != test.expected {
			t.Errorf("cleanManagerPrefix(%s, %s) = %s; want %s",
				test.path, test.managerName, result, test.expected)
		}
	}
}

// Test helpers and setup functions
var (
	originalMu          sync.Mutex
	originalComposites  []*Folder
	originalObjectMap   map[string]map[string]object
	originalManagerPath string
)

func setupTest(t *testing.T) {
	projectRoot := findProjectRoot(t)
	os.Chdir(projectRoot)

	originalMu = mu
	originalComposites = Composites
	originalObjectMap = ObjectMap
	originalManagerPath = managersFilePath

	mu = sync.Mutex{}
	Composites = []*Folder{}
	ObjectMap = make(map[string]map[string]object)
	managersFilePath = "test_managers.json"
}

func cleanupTest(t *testing.T) {
	mu = originalMu
	Composites = originalComposites
	ObjectMap = originalObjectMap
	managersFilePath = originalManagerPath

	os.Remove("test_managers.json")
}

func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if filepath.Base(dir) == "Smart-File-Manager" {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root")
		}
		dir = parent
	}
}
