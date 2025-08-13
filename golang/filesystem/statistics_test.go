package filesystem

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

// Mock data structures for testing
type MockFolder struct {
	Name       string
	Files      []*MockFile
	Subfolders []*MockFolder
}

type MockFile struct {
	Path string
}

// Test helper functions
func createTempFile(t *testing.T, dir, name string, content []byte, modTime time.Time) string {
	filePath := filepath.Join(dir, name)
	err := ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filePath, err)
	}

	// Set modification time
	err = os.Chtimes(filePath, modTime, modTime)
	if err != nil {
		t.Fatalf("Failed to set modification time for %s: %v", filePath, err)
	}

	return filePath
}

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "filesystem_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return dir
}

func setupTestFiles(t *testing.T, dir string) []string {
	now := time.Now()

	files := []string{
		createTempFile(t, dir, "recent1.txt", []byte("small file"), now.Add(-1*time.Hour)),
		createTempFile(t, dir, "recent2.txt", []byte("medium content here"), now.Add(-2*time.Hour)),
		createTempFile(t, dir, "old1.txt", []byte("old"), now.Add(-24*time.Hour)),
		createTempFile(t, dir, "old2.txt", []byte("very old"), now.Add(-48*time.Hour)),
		createTempFile(t, dir, "large.txt", make([]byte, 1000), now.Add(-3*time.Hour)),
	}

	return files
}

// TestCalculateTotalSize tests the calculateTotalSize function
func TestCalculateTotalSize(t *testing.T) {
	tests := []struct {
		name     string
		files    []fileInfo
		expected int64
	}{
		{
			name:     "empty files",
			files:    []fileInfo{},
			expected: 0,
		},
		{
			name: "single file",
			files: []fileInfo{
				{size: 100},
			},
			expected: 100,
		},
		{
			name: "multiple files",
			files: []fileInfo{
				{size: 100},
				{size: 200},
				{size: 300},
			},
			expected: 600,
		},
		{
			name: "files with zero size",
			files: []fileInfo{
				{size: 0},
				{size: 100},
				{size: 0},
			},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateTotalSize(tt.files)
			if result != tt.expected {
				t.Errorf("calculateTotalSize() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

// TestCalculateUmbrellaCounts tests the calculateUmbrellaCounts function
func TestCalculateUmbrellaCounts(t *testing.T) {
	tests := []struct {
		name     string
		files    []fileInfo
		expected []int
	}{
		{
			name:     "empty files",
			files:    []fileInfo{},
			expected: []int{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "all document types",
			files: []fileInfo{
				{umbrella: "Documents"},
				{umbrella: "Images"},
				{umbrella: "Music"},
				{umbrella: "Presentations"},
				{umbrella: "Videos"},
				{umbrella: "Spreadsheets"},
				{umbrella: "Archives"},
				{umbrella: "Unknown"},
			},
			expected: []int{1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name: "multiple of same type",
			files: []fileInfo{
				{umbrella: "Documents"},
				{umbrella: "Documents"},
				{umbrella: "Images"},
				{umbrella: "Unknown"},
				{umbrella: "SomeOtherType"}, // Should count as Unknown
			},
			expected: []int{2, 1, 0, 0, 0, 0, 0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateUmbrellaCounts(tt.files)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("calculateUmbrellaCounts() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCountFolders tests the countFolders function
func TestCountFolders(t *testing.T) {
	tests := []struct {
		name     string
		folder   *Folder
		expected int
	}{
		{
			name: "no subfolders",
			folder: &Folder{
				Subfolders: []*Folder{},
			},
			expected: 0,
		},
		{
			name: "single level subfolders",
			folder: &Folder{
				Subfolders: []*Folder{
					{},
					{},
					{},
				},
			},
			expected: 3,
		},
		{
			name: "nested subfolders",
			folder: &Folder{
				Subfolders: []*Folder{
					{
						Subfolders: []*Folder{
							{},
							{},
						},
					},
					{
						Subfolders: []*Folder{
							{},
						},
					},
				},
			},
			expected: 5, // 2 direct + 3 nested
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countFolders(tt.folder)
			if result != tt.expected {
				t.Errorf("countFolders() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

// TestGetNewestFiles tests the getNewestFiles function
func TestGetNewestFiles(t *testing.T) {
	now := time.Now()
	files := []fileInfo{
		{path: "path1", name: "file1.txt", modTime: now.Add(-1 * time.Hour)},
		{path: "path2", name: "file2.txt", modTime: now.Add(-2 * time.Hour)},
		{path: "path3", name: "file3.txt", modTime: now.Add(-3 * time.Hour)},
		{path: "path4", name: "file4.txt", modTime: now.Add(-4 * time.Hour)},
	}

	tests := []struct {
		name     string
		files    []fileInfo
		limit    int
		expected []file
	}{
		{
			name:     "empty files",
			files:    []fileInfo{},
			limit:    5,
			expected: []file{},
		},
		{
			name:  "limit less than file count",
			files: files,
			limit: 2,
			expected: []file{
				{FilePath: "path1", FileName: "file1.txt"},
				{FilePath: "path2", FileName: "file2.txt"},
			},
		},
		{
			name:  "limit greater than file count",
			files: files[:2],
			limit: 5,
			expected: []file{
				{FilePath: "path1", FileName: "file1.txt"},
				{FilePath: "path2", FileName: "file2.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNewestFiles(tt.files, tt.limit)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getNewestFiles() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetOldestFiles tests the getOldestFiles function
// func TestGetOldestFiles(t *testing.T) {
// 	now := time.Now()
// 	files := []fileInfo{
// 		{path: "path1", name: "file1.txt", modTime: now.Add(-1 * time.Hour)},
// 		{path: "path2", name: "file2.txt", modTime: now.Add(-2 * time.Hour)},
// 		{path: "path3", name: "file3.txt", modTime: now.Add(-3 * time.Hour)},
// 		{path: "path4", name: "file4.txt", modTime: now.Add(-4 * time.Hour)},
// 	}

// 	tests := []struct {
// 		name     string
// 		files    []fileInfo
// 		limit    int
// 		expected []file
// 	}{
// 		{
// 			name:     "empty files",
// 			files:    []fileInfo{},
// 			limit:    5,
// 			expected: []file{},
// 		},
// 		{
// 			name:  "limit less than file count",
// 			files: files,
// 			limit: 2,
// 			expected: []file{
// 				{FilePath: "path4", FileName: "file4.txt"},
// 				{FilePath: "path3", FileName: "file3.txt"},
// 			},
// 		},
// 		{
// 			name:  "limit greater than file count",
// 			files: files[:2],
// 			limit: 5,
// 			expected: []file{
// 				{FilePath: "path2", FileName: "file2.txt"},
// 				{FilePath: "path1", FileName: "file1.txt"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result := getOldestFiles(tt.files, tt.limit)
// 			if !reflect.DeepEqual(result, tt.expected) {
// 				t.Errorf("getOldestFiles() = %v, expected %v", result, tt.expected)
// 			}
// 		})
// 	}
// }

// TestGetLargestFiles tests the getLargestFiles function
func TestGetLargestFiles(t *testing.T) {
	files := []fileInfo{
		{path: "path1", name: "file1.txt", size: 100},
		{path: "path2", name: "file2.txt", size: 200},
		{path: "path3", name: "file3.txt", size: 300},
		{path: "path4", name: "file4.txt", size: 400},
	}

	tests := []struct {
		name     string
		files    []fileInfo
		limit    int
		expected []file
	}{
		{
			name:     "empty files",
			files:    []fileInfo{},
			limit:    5,
			expected: []file{},
		},
		{
			name:  "limit less than file count",
			files: files,
			limit: 2,
			expected: []file{
				{FilePath: "path4", FileName: "file4.txt"},
				{FilePath: "path3", FileName: "file3.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLargestFiles(tt.files, tt.limit)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getLargestFiles() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCollectFilesRecursive tests the collectFilesRecursive function
func TestCollectFilesRecursive(t *testing.T) {
	// Create temporary directory and files
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	now := time.Now()
	file1 := createTempFile(t, tempDir, "file1.txt", []byte("content1"), now)
	file2 := createTempFile(t, subDir, "file2.txt", []byte("content2"), now)

	// Create folder structure
	folder := &Folder{
		Name: "testfolder",
		Files: []*File{
			{Path: file1},
		},
		Subfolders: []*Folder{
			{
				Name: "subfolder",
				Files: []*File{
					{Path: file2},
				},
				Subfolders: []*Folder{},
			},
		},
	}

	var files []fileInfo
	collectFilesRecursive(folder, &files)

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Sort files by path for consistent testing
	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})

	expectedFiles := []string{file1, file2}
	sort.Strings(expectedFiles)

	for i, file := range files {
		if file.path != expectedFiles[i] {
			t.Errorf("Expected file path %s, got %s", expectedFiles[i], file.path)
		}
		if file.name != filepath.Base(expectedFiles[i]) {
			t.Errorf("Expected file name %s, got %s", filepath.Base(expectedFiles[i]), file.name)
		}
	}
}

// Mock interfaces for dependency injection
type CompositeProvider interface {
	GetComposites() []*Folder
}

type TypeLoader interface {
	LoadTypes(folder *Folder, name string)
}

// Mock implementations
type MockCompositeProvider struct {
	folders []*Folder
}

func (m *MockCompositeProvider) GetComposites() []*Folder {
	return m.folders
}

type MockTypeLoader struct{}

func (m *MockTypeLoader) LoadTypes(folder *Folder, name string) {
	// Mock implementation - do nothing
}

// Modified version of collectManagerFiles for testing
func collectManagerFilesWithMocks(folder *Folder, typeLoader TypeLoader) []fileInfo {
	var files []fileInfo

	log.Printf("LoadTypes: Starting for folder %s", folder.Name)
	typeLoader.LoadTypes(folder, folder.Name)
	log.Printf("LoadTypes: Completed for folder %s", folder.Name)

	log.Printf("collectFilesRecursive: Starting for folder %s", folder.Name)
	collectFilesRecursive(folder, &files)
	log.Printf("collectFilesRecursive: Completed for folder %s, found %d files", folder.Name, len(files))

	return files
}

// TestStatHandlerIntegration tests the core logic without HTTP
func TestStatHandlerLogic(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	file1 := createTempFile(t, tempDir, "test1.txt", []byte("test content"), time.Now())

	// Create mock data
	mockProvider := &MockCompositeProvider{
		folders: []*Folder{
			{
				Name: "TestManager",
				Files: []*File{
					{Path: file1},
				},
				Subfolders: []*Folder{},
			},
		},
	}

	mockTypeLoader := &MockTypeLoader{}

	// Test the core logic
	composites := mockProvider.GetComposites()
	if len(composites) != 1 {
		t.Fatalf("Expected 1 composite, got %d", len(composites))
	}

	folder := composites[0]
	manager := ManagerStatistics{
		ManagerName: folder.Name,
	}

	// Collect all file information for this manager
	allFiles := collectManagerFilesWithMocks(folder, mockTypeLoader)
	if len(allFiles) != 1 {
		t.Errorf("Expected 1 file, got %d", len(allFiles))
	}

	// Calculate statistics
	manager.Files = len(allFiles)
	manager.Folders = countFolders(folder)
	manager.Size = calculateTotalSize(allFiles)
	manager.UmbrellaCounts = calculateUmbrellaCounts(allFiles)

	// Get file rankings
	manager.Recent = getNewestFiles(allFiles, 5)
	manager.Oldest = getOldestFiles(allFiles, 5)
	manager.Largest = getLargestFiles(allFiles, 5)

	// Verify results
	if manager.ManagerName != "TestManager" {
		t.Errorf("Expected manager name 'TestManager', got '%s'", manager.ManagerName)
	}

	if manager.Files != 1 {
		t.Errorf("Expected 1 file, got %d", manager.Files)
	}

	if manager.Folders != 0 {
		t.Errorf("Expected 0 folders, got %d", manager.Folders)
	}

	if len(manager.Recent) != 1 {
		t.Errorf("Expected 1 recent file, got %d", len(manager.Recent))
	}
}

// TestStatHandlerHTTP tests the HTTP response format (without mocking internal functions)
func TestStatHandlerHTTP(t *testing.T) {
	// Skip this test if dependencies are not available
	// This test would need to be run in an environment where GetComposites() works
	t.Skip("Skipping HTTP integration test - requires real dependencies")

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StatHandler)

	handler.ServeHTTP(rr, req)

	// Check that we get a valid JSON response
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", rr.Header().Get("Content-Type"))
	}

	// Try to unmarshal the response
	var managers []ManagerStatistics
	err = json.Unmarshal(rr.Body.Bytes(), &managers)
	if err != nil {
		t.Errorf("Failed to unmarshal response as JSON: %v", err)
	}
}

// Benchmark tests
func BenchmarkCalculateTotalSize(b *testing.B) {
	files := make([]fileInfo, 1000)
	for i := range files {
		files[i] = fileInfo{size: int64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateTotalSize(files)
	}
}

func BenchmarkCalculateUmbrellaCounts(b *testing.B) {
	files := make([]fileInfo, 1000)
	umbrellaTypes := []string{"Documents", "Images", "Music", "Videos", "Unknown"}

	for i := range files {
		files[i] = fileInfo{umbrella: umbrellaTypes[i%len(umbrellaTypes)]}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateUmbrellaCounts(files)
	}
}

func BenchmarkGetNewestFiles(b *testing.B) {
	files := make([]fileInfo, 1000)
	now := time.Now()

	for i := range files {
		files[i] = fileInfo{
			path:    "path" + string(rune(i)),
			name:    "file" + string(rune(i)),
			modTime: now.Add(time.Duration(i) * time.Hour),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getNewestFiles(files, 10)
	}
}

// Table-driven test for file ranking functions
func TestFileRankingFunctions(t *testing.T) {
	now := time.Now()
	testFiles := []fileInfo{
		{path: "path1", name: "file1.txt", size: 100, modTime: now.Add(-1 * time.Hour)},
		{path: "path2", name: "file2.txt", size: 300, modTime: now.Add(-3 * time.Hour)},
		{path: "path3", name: "file3.txt", size: 200, modTime: now.Add(-2 * time.Hour)},
	}

	tests := []struct {
		name     string
		function func([]fileInfo, int) []file
		expected []string // Expected file names in order
	}{
		{
			name:     "newest files",
			function: getNewestFiles,
			expected: []string{"file1.txt", "file3.txt", "file2.txt"},
		},
		{
			name:     "oldest files",
			function: getOldestFiles,
			expected: []string{"file2.txt", "file3.txt", "file1.txt"},
		},
		{
			name:     "largest files",
			function: getLargestFiles,
			expected: []string{"file2.txt", "file3.txt", "file1.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(testFiles, 3)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d files, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i].FileName != expected {
					t.Errorf("Position %d: expected %s, got %s", i, expected, result[i].FileName)
				}
			}
		})
	}
}
