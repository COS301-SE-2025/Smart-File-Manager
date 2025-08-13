package filesystem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Test helper functions
func createTestFolder() *Folder {

	files := []*File{
		{
			Name: "report.pdf",
			Path: "/home/user/documents/report.pdf",
		},
		{
			Name: "vacation.jpg",
			Path: "/home/user/photos/vacation.jpg",
		},
		{
			Name: "song.mp3",
			Path: "/home/user/music/song.mp3",
		},
	}
	folder := &Folder{
		Name:  "test-folder",
		Files: files,
	}

	return folder
}

func createTestBulkList() []TagsStruct {
	return []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"work", "important", "pdf"},
		},
		{
			FilePath: "/home/user/photos/vacation.jpg",
			Tags:     []string{"holiday", "family", "2025"},
		},
		{
			FilePath: "/home/user/music/song.mp3",
			Tags:     []string{"music", "mp3", "favorites"},
		},
	}
}

// Tests for BulkAddTags
func TestBulkAddTags_Success(t *testing.T) {
	folder := createTestFolder()
	bulkList := createTestBulkList()

	err := BulkAddTags(folder, bulkList)

	if err != nil {
		t.Errorf("BulkAddTags returned error: %v", err)
	}

	// Verify tags were added correctly
	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil {
		if !contains(pdfFile.Tags, "work") || !contains(pdfFile.Tags, "important") || !contains(pdfFile.Tags, "pdf") {
			t.Errorf("PDF file missing expected tags. Current tags: %v", pdfFile.Tags)
		}
	}

	jpgFile := folder.GetFile("/home/user/photos/vacation.jpg")
	if jpgFile != nil {
		if !contains(jpgFile.Tags, "holiday") || !contains(jpgFile.Tags, "family") || !contains(jpgFile.Tags, "2025") {
			t.Errorf("JPG file missing expected tags. Current tags: %v", jpgFile.Tags)
		}
	}

	mp3File := folder.GetFile("/home/user/music/song.mp3")
	if mp3File != nil {
		if !contains(mp3File.Tags, "music") || !contains(mp3File.Tags, "mp3") || !contains(mp3File.Tags, "favorites") {
			t.Errorf("MP3 file missing expected tags. Current tags: %v", mp3File.Tags)
		}
	}
}

func TestBulkAddTags_EmptyBulkList(t *testing.T) {
	folder := createTestFolder()
	emptyBulkList := []TagsStruct{}

	err := BulkAddTags(folder, emptyBulkList)

	if err != nil {
		t.Errorf("BulkAddTags returned error for empty list: %v", err)
	}

	// Verify no tags were added - this test assumes files start with no tags
	// If your createTestFolder() adds files with existing tags, adjust accordingly
	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil && len(pdfFile.Tags) != 0 {
		t.Errorf("Expected no tags on PDF file, but found: %v", pdfFile.Tags)
	}

	jpgFile := folder.GetFile("/home/user/photos/vacation.jpg")
	if jpgFile != nil && len(jpgFile.Tags) != 0 {
		t.Errorf("Expected no tags on JPG file, but found: %v", jpgFile.Tags)
	}

	mp3File := folder.GetFile("/home/user/music/song.mp3")
	if mp3File != nil && len(mp3File.Tags) != 0 {
		t.Errorf("Expected no tags on MP3 file, but found: %v", mp3File.Tags)
	}
}

func TestBulkAddTags_DuplicateTags(t *testing.T) {
	folder := createTestFolder()

	// Pre-add some tags
	folder.AddTagToFile("/home/user/documents/report.pdf", "work")

	bulkList := []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"work", "important", "pdf"}, // "work" is duplicate
		},
	}

	err := BulkAddTags(folder, bulkList)

	if err != nil {
		t.Errorf("BulkAddTags returned error: %v", err)
	}

	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil {
		// Count occurrences of "work" tag
		workCount := 0
		for _, tag := range pdfFile.Tags {
			if tag == "work" {
				workCount++
			}
		}

		if workCount != 1 {
			t.Errorf("Expected 'work' tag to appear once, but found %d occurrences. Tags: %v", workCount, pdfFile.Tags)
		}

		// Verify other tags were added
		if !contains(pdfFile.Tags, "important") || !contains(pdfFile.Tags, "pdf") {
			t.Errorf("Expected tags missing from PDF file. Current tags: %v", pdfFile.Tags)
		}
	}
}

func TestBulkAddTags_NonExistentFile(t *testing.T) {
	folder := createTestFolder()
	bulkList := []TagsStruct{
		{
			FilePath: "/nonexistent/file.txt",
			Tags:     []string{"test"},
		},
	}

	// This should not return an error based on current implementation
	// The AddTagToFile method handles non-existent files internally
	err := BulkAddTags(folder, bulkList)

	// Note: Current implementation doesn't check for file existence in BulkAddTags
	// If you want it to fail on non-existent files, modify the implementation
	if err != nil {
		t.Errorf("BulkAddTags returned unexpected error: %v", err)
	}
}

func TestBulkAddTags_NilFolder(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when passing nil folder, but didn't panic")
		}
	}()

	bulkList := createTestBulkList()
	BulkAddTags(nil, bulkList)
}

// Tests for BulkRemoveTags
func TestBulkRemoveTags_Success(t *testing.T) {
	folder := createTestFolder()

	// First add tags
	bulkAddList := createTestBulkList()
	BulkAddTags(folder, bulkAddList)

	// Then remove some tags
	bulkRemoveList := []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"work", "important"},
		},
		{
			FilePath: "/home/user/photos/vacation.jpg",
			Tags:     []string{"holiday"},
		},
	}

	err := BulkRemoveTags(folder, bulkRemoveList)

	if err != nil {
		t.Errorf("BulkRemoveTags returned error: %v", err)
	}

	// Verify tags were removed
	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil {
		if contains(pdfFile.Tags, "work") || contains(pdfFile.Tags, "important") {
			t.Errorf("Tags were not removed from PDF file. Current tags: %v", pdfFile.Tags)
		}
		if !contains(pdfFile.Tags, "pdf") {
			t.Errorf("Expected 'pdf' tag to remain on PDF file. Current tags: %v", pdfFile.Tags)
		}
	}

	jpgFile := folder.GetFile("/home/user/photos/vacation.jpg")
	if jpgFile != nil {
		if contains(jpgFile.Tags, "holiday") {
			t.Errorf("'holiday' tag was not removed from JPG file. Current tags: %v", jpgFile.Tags)
		}
		if !contains(jpgFile.Tags, "family") || !contains(jpgFile.Tags, "2025") {
			t.Errorf("Expected remaining tags on JPG file. Current tags: %v", jpgFile.Tags)
		}
	}
}

func TestBulkRemoveTags_NonExistentFile(t *testing.T) {
	folder := createTestFolder()
	bulkList := []TagsStruct{
		{
			FilePath: "/nonexistent/file.txt",
			Tags:     []string{"test"},
		},
	}

	err := BulkRemoveTags(folder, bulkList)

	if err == nil {
		t.Errorf("Expected error for non-existent file, but got nil")
	}

	expectedError := "file not found: /nonexistent/file.txt"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestBulkRemoveTags_NonExistentTag(t *testing.T) {
	folder := createTestFolder()

	// Add a file with some tags
	folder.AddTagToFile("/home/user/documents/report.pdf", "existing-tag")

	bulkList := []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"non-existent-tag"},
		},
	}

	err := BulkRemoveTags(folder, bulkList)

	if err == nil {
		t.Errorf("Expected error for non-existent tag, but got nil")
	}

	expectedError := "failed to remove tag non-existent-tag from file /home/user/documents/report.pdf"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestBulkRemoveTags_EmptyBulkList(t *testing.T) {
	folder := createTestFolder()
	emptyBulkList := []TagsStruct{}

	err := BulkRemoveTags(folder, emptyBulkList)

	if err != nil {
		t.Errorf("BulkRemoveTags returned error for empty list: %v", err)
	}
}

func TestBulkRemoveTags_PartialFailure(t *testing.T) {
	folder := createTestFolder()

	// Add some tags first
	folder.AddTagToFile("/home/user/documents/report.pdf", "work")

	bulkList := []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"work", "non-existent-tag"}, // First succeeds, second fails
		},
	}

	err := BulkRemoveTags(folder, bulkList)

	if err == nil {
		t.Errorf("Expected error for non-existent tag, but got nil")
	}

	// Verify that the function stops on first failure
	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil {
		if contains(pdfFile.Tags, "work") {
			t.Errorf("Expected 'work' tag to be removed before failure. Current tags: %v", pdfFile.Tags)
		}
	}
}

func TestBulkRemoveTags_NilFolder(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when passing nil folder, but didn't panic")
		}
	}()

	bulkList := createTestBulkList()
	BulkRemoveTags(nil, bulkList)
}

// Integration tests
func TestBulkOperations_Integration(t *testing.T) {
	folder := createTestFolder()

	// Test add then remove cycle
	addList := []TagsStruct{
		{
			FilePath: "/home/user/documents/report.pdf",
			Tags:     []string{"work", "important", "pdf", "draft"},
		},
	}

	// Add tags
	err := BulkAddTags(folder, addList)
	if err != nil {
		t.Errorf("Failed to add tags: %v", err)
	}

	// Verify all tags were added
	pdfFile := folder.GetFile("/home/user/documents/report.pdf")
	if pdfFile != nil {
		expectedTags := []string{"work", "important", "pdf", "draft"}
		for _, expectedTag := range expectedTags {
			if !contains(pdfFile.Tags, expectedTag) {
				t.Errorf("Expected tag '%s' not found. Current tags: %v", expectedTag, pdfFile.Tags)
			}
		}

		// Remove some tags
		removeList := []TagsStruct{
			{
				FilePath: "/home/user/documents/report.pdf",
				Tags:     []string{"draft", "important"},
			},
		}

		err = BulkRemoveTags(folder, removeList)
		if err != nil {
			t.Errorf("Failed to remove tags: %v", err)
		}

		// Verify correct tags remain
		if contains(pdfFile.Tags, "draft") || contains(pdfFile.Tags, "important") {
			t.Errorf("Removed tags still present. Current tags: %v", pdfFile.Tags)
		}
		if !contains(pdfFile.Tags, "work") || !contains(pdfFile.Tags, "pdf") {
			t.Errorf("Expected tags missing after removal. Current tags: %v", pdfFile.Tags)
		}
	}
}

// Benchmark tests
func BenchmarkBulkAddTags(b *testing.B) {
	// folder := createTestFolder()
	bulkList := createTestBulkList()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create fresh folder for each iteration to avoid tag duplication effects
		testFolder := createTestFolder()
		BulkAddTags(testFolder, bulkList)
	}
}

func BenchmarkBulkRemoveTags(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		folder := createTestFolder()
		bulkList := createTestBulkList()
		BulkAddTags(folder, bulkList) // Setup tags to remove
		b.StartTimer()

		BulkRemoveTags(folder, bulkList)
	}
}

func TestAPI_BulkDeleteFileHandler_MixedExistentNonExistent(t *testing.T) {
	// Setup: Create temp folder and one test file
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create only one file
	existingFile := filepath.Join(dataDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Register the file in composites
	testFolder := &Folder{
		Name:    "mixedTest",
		Path:    dataDir,
		NewPath: dataDir,
		Files: []*File{
			{Name: "existing.txt", Path: existingFile},
		},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Create JSON payload with both existing and non-existing files
	nonExistentFile := filepath.Join(dataDir, "nonexistent.txt")
	jsonBody := `[
		{"file_path": "` + existingFile + `"},
		{"file_path": "` + nonExistentFile + `"}
	]`

	req := httptest.NewRequest("POST", "/bulkDeleteFile?name=mixedTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFileHandler(w, req)

	// The handler should handle this gracefully - os.RemoveAll won't fail on non-existent files
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify existing file was deleted
	if _, err := os.Stat(existingFile); !os.IsNotExist(err) {
		t.Error("existing file should have been deleted")
	}
}

func TestAPI_BulkDeleteFolderHandler_MixedExistentNonExistent(t *testing.T) {
	// Setup: Create temp folder and one test folder
	tmp := t.TempDir()
	projectRoot := filepath.Join(tmp, "Smart-File-Manager")
	dataDir := filepath.Join(projectRoot, "data")
	existingFolder := filepath.Join(dataDir, "existing")

	if err := os.MkdirAll(existingFolder, 0755); err != nil {
		t.Fatal(err)
	}

	// Register the folder in composites
	testFolder := &Folder{
		Name:    "mixedFolderTest",
		Path:    dataDir,
		NewPath: dataDir,
		Subfolders: []*Folder{
			{Name: "existing", Path: existingFolder},
		},
	}

	Composites = []*Folder{testFolder}

	// Change working directory
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Create JSON payload with both existing and non-existing folders
	nonExistentFolder := filepath.Join(dataDir, "nonexistent")
	jsonBody := `[
		{"file_path": "` + existingFolder + `"},
		{"file_path": "` + nonExistentFolder + `"}
	]`

	req := httptest.NewRequest("POST", "/bulkDeleteFolder?name=mixedFolderTest", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	BulkDeleteFolderHandler(w, req)

	// The handler should handle this gracefully - os.RemoveAll won't fail on non-existent folders
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify existing folder was deleted
	if _, err := os.Stat(existingFolder); !os.IsNotExist(err) {
		t.Error("existing folder should have been deleted")
	}
}

// Test helper function to create a test folder with various file types
func createTestFolderWithTypes() *Folder {
	files := []*File{
		{Name: "document.pdf", Path: "/test/documents/document.pdf"},
		{Name: "report.docx", Path: "/test/documents/report.docx"},
		{Name: "photo.jpg", Path: "/test/images/photo.jpg"},
		{Name: "logo.png", Path: "/test/images/logo.png"},
		{Name: "song.mp3", Path: "/test/music/song.mp3"},
		{Name: "audio.wav", Path: "/test/music/audio.wav"},
		{Name: "presentation.pptx", Path: "/test/presentations/presentation.pptx"},
		{Name: "video.mp4", Path: "/test/videos/video.mp4"},
		{Name: "spreadsheet.xlsx", Path: "/test/spreadsheets/spreadsheet.xlsx"},
		{Name: "archive.zip", Path: "/test/archives/archive.zip"},
		{Name: "unknown.xyz", Path: "/test/unknown/unknown.xyz"},
	}

	subfolder := &Folder{
		Name: "subfolder",
		Path: "/test/subfolder",
		Files: []*File{
			{Name: "nested.txt", Path: "/test/subfolder/nested.txt"},
			{Name: "nested_image.gif", Path: "/test/subfolder/nested_image.gif"},
		},
	}

	folder := &Folder{
		Name:       "test-folder-types",
		Path:       "/test",
		Files:      files,
		Subfolders: []*Folder{subfolder},
	}

	return folder
}

// Tests for GetCategory function
func TestGetCategory(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		// Documents
		{"document.pdf", "Documents"},
		{"report.docx", "Documents"},
		{"text.txt", "Documents"},
		{"spreadsheet.csv", "Documents"},
		{"book.epub", "Documents"},

		// Images
		{"photo.jpg", "Images"},
		{"image.jpeg", "Images"},
		{"logo.png", "Images"},
		{"animation.gif", "Images"},
		{"vector.svg", "Images"},
		{"raw.cr2", "Images"},

		// Music
		{"song.mp3", "Music"},
		{"audio.wav", "Music"},
		{"lossless.flac", "Music"},
		{"compressed.aac", "Music"},
		{"midi.mid", "Music"},

		// Videos
		{"movie.mp4", "Videos"},
		{"clip.avi", "Videos"},
		{"web.webm", "Videos"},
		{"mobile.3gp", "Videos"},

		// Presentations
		{"slides.pptx", "Presentations"},
		{"keynote.key", "Presentations"},
		{"presentation.odp", "Presentations"},

		// Spreadsheets
		{"data.xlsx", "Spreadsheets"},
		{"old.xls", "Spreadsheets"},
		{"calc.ods", "Spreadsheets"},
		{"values.tsv", "Spreadsheets"},

		// Archives
		{"compressed.zip", "Archives"},
		{"backup.tar", "Archives"},
		{"compressed.7z", "Archives"},
		{"disk.iso", "Archives"},

		// Unknown/Edge cases
		{"unknown.xyz", "Unknown"},
		{"file.unknown", "Unknown"},
		{"no_extension", "Unknown"},
		{"", "Unknown"},
		{".hidden", "Unknown"},

		// Case sensitivity tests
		{"FILE.PDF", "Documents"},
		{"IMAGE.JPG", "Images"},
		{"SONG.MP3", "Music"},

		// Multiple extensions
		{"archive.tar.gz", "Archives"}, // Should get "gz" extension
		{"backup.tar.bz2", "Archives"}, // Should get "bz2" extension

		// Path with filename
		{"/path/to/document.pdf", "Documents"},
		{"/home/user/music/song.mp3", "Music"},
		{"C:\\Windows\\file.txt", "Documents"},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := GetCategory(test.filename)
			if result != test.expected {
				t.Errorf("GetCategory(%q) = %q, expected %q", test.filename, result, test.expected)
			}
		})
	}
}

// Tests for LoadTypes function with map of maps
func TestLoadTypes_EmptyFolder(t *testing.T) {
	// Clear and initialize ObjectMap
	ObjectMap = make(map[string]map[string]object)

	emptyFolder := &Folder{
		Name:       "empty-manager",
		Files:      []*File{},
		Subfolders: []*Folder{},
	}

	// Initialize the inner map for this manager
	ObjectMap["empty-manager"] = make(map[string]object)

	LoadTypes(emptyFolder, "empty-manager")

	if len(ObjectMap["empty-manager"]) != 0 {
		t.Errorf("Expected empty ObjectMap for manager 'empty-manager', but got %d entries", len(ObjectMap["empty-manager"]))
	}
}

func TestLoadTypes_FolderWithFiles(t *testing.T) {
	// Clear and initialize ObjectMap
	ObjectMap = make(map[string]map[string]object)

	folder := createTestFolderWithTypes()
	managerName := "test-manager"

	// Initialize the inner map for this manager
	ObjectMap[managerName] = make(map[string]object)

	LoadTypes(folder, managerName)

	// Check that all files were loaded into ObjectMap
	expectedEntries := 13 // 11 files in main folder + 2 files in subfolder
	if len(ObjectMap[managerName]) != expectedEntries {
		t.Errorf("Expected %d entries in ObjectMap for manager %q, but got %d", expectedEntries, managerName, len(ObjectMap[managerName]))
	}

	// Test specific file entries
	tests := []struct {
		path             string
		expectedFileType string
		expectedUmbrella string
	}{
		{"/test/documents/document.pdf", "pdf", "Documents"},
		{"/test/images/photo.jpg", "jpg", "Images"},
		{"/test/music/song.mp3", "mp3", "Music"},
		{"/test/videos/video.mp4", "mp4", "Videos"},
		{"/test/archives/archive.zip", "zip", "Archives"},
		{"/test/unknown/unknown.xyz", "xyz", "Unknown"},
		{"/test/subfolder/nested.txt", "txt", "Documents"},
		{"/test/subfolder/nested_image.gif", "gif", "Images"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			obj, exists := ObjectMap[managerName][test.path]
			if !exists {
				t.Errorf("Expected path %q to exist in ObjectMap for manager %q", test.path, managerName)
				return
			}

			if obj.fileType != test.expectedFileType {
				t.Errorf("Expected fileType %q, got %q", test.expectedFileType, obj.fileType)
			}

			if obj.umbrellaType != test.expectedUmbrella {
				t.Errorf("Expected umbrellaType %q, got %q", test.expectedUmbrella, obj.umbrellaType)
			}
		})
	}
}

func TestLoadTypes_MultipleManagers(t *testing.T) {
	// Clear and initialize ObjectMap
	ObjectMap = make(map[string]map[string]object)

	// Create two different folders/managers
	folder1 := &Folder{
		Name: "manager1",
		Files: []*File{
			{Name: "doc1.pdf", Path: "/manager1/doc1.pdf"},
			{Name: "image1.jpg", Path: "/manager1/image1.jpg"},
		},
	}

	folder2 := &Folder{
		Name: "manager2",
		Files: []*File{
			{Name: "doc2.txt", Path: "/manager2/doc2.txt"},
			{Name: "music1.mp3", Path: "/manager2/music1.mp3"},
		},
	}

	// Initialize inner maps
	ObjectMap["manager1"] = make(map[string]object)
	ObjectMap["manager2"] = make(map[string]object)

	// Load types for both managers
	LoadTypes(folder1, "manager1")
	LoadTypes(folder2, "manager2")

	// Verify both managers have their own separate maps
	if len(ObjectMap["manager1"]) != 2 {
		t.Errorf("Expected 2 entries for manager1, got %d", len(ObjectMap["manager1"]))
	}

	if len(ObjectMap["manager2"]) != 2 {
		t.Errorf("Expected 2 entries for manager2, got %d", len(ObjectMap["manager2"]))
	}

	// Verify manager1 files
	if obj, exists := ObjectMap["manager1"]["/manager1/doc1.pdf"]; exists {
		if obj.fileType != "pdf" || obj.umbrellaType != "Documents" {
			t.Errorf("Manager1 PDF file incorrect: fileType=%q, umbrellaType=%q", obj.fileType, obj.umbrellaType)
		}
	} else {
		t.Error("Manager1 PDF file not found")
	}

	// Verify manager2 files
	if obj, exists := ObjectMap["manager2"]["/manager2/music1.mp3"]; exists {
		if obj.fileType != "mp3" || obj.umbrellaType != "Music" {
			t.Errorf("Manager2 MP3 file incorrect: fileType=%q, umbrellaType=%q", obj.fileType, obj.umbrellaType)
		}
	} else {
		t.Error("Manager2 MP3 file not found")
	}

	// Verify separation - manager1 shouldn't have manager2's files
	if _, exists := ObjectMap["manager1"]["/manager2/music1.mp3"]; exists {
		t.Error("Manager1 should not contain manager2's files")
	}
}

func TestLoadTypes_FileWithoutExtension(t *testing.T) {
	// Clear and initialize ObjectMap
	ObjectMap = make(map[string]map[string]object)

	folder := &Folder{
		Name: "test-manager",
		Files: []*File{
			{Name: "no_extension", Path: "/test/no_extension"},
		},
	}

	// Initialize inner map
	ObjectMap["test-manager"] = make(map[string]object)

	// This should cause a panic because strings.Split will not have index [1]
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when processing file without extension, but didn't panic")
		} else {
			t.Logf("Function correctly panics on files without extension: %v", r)
		}
	}()

	LoadTypes(folder, "test-manager")
}

func TestLoadTypes_NestedFolders(t *testing.T) {
	// Clear and initialize ObjectMap
	ObjectMap = make(map[string]map[string]object)

	// Create deeply nested folder structure
	deepNested := &Folder{
		Name: "deep",
		Files: []*File{
			{Name: "deep.txt", Path: "/test/level1/level2/deep.txt"},
		},
	}

	level2 := &Folder{
		Name: "level2",
		Files: []*File{
			{Name: "level2.jpg", Path: "/test/level1/level2.jpg"},
		},
		Subfolders: []*Folder{deepNested},
	}

	level1 := &Folder{
		Name: "level1",
		Files: []*File{
			{Name: "level1.pdf", Path: "/test/level1.pdf"},
		},
		Subfolders: []*Folder{level2},
	}

	rootFolder := &Folder{
		Name: "root",
		Files: []*File{
			{Name: "root.mp3", Path: "/test/root.mp3"},
		},
		Subfolders: []*Folder{level1},
	}

	managerName := "nested-manager"
	ObjectMap[managerName] = make(map[string]object)

	LoadTypes(rootFolder, managerName)

	// Should have 4 files total
	if len(ObjectMap[managerName]) != 4 {
		t.Errorf("Expected 4 entries in ObjectMap for manager %q, but got %d", managerName, len(ObjectMap[managerName]))
	}

	// Verify all files are present
	expectedPaths := []string{
		"/test/root.mp3",
		"/test/level1.pdf",
		"/test/level1/level2.jpg",
		"/test/level1/level2/deep.txt",
	}

	for _, path := range expectedPaths {
		if _, exists := ObjectMap[managerName][path]; !exists {
			t.Errorf("Expected path %q to exist in ObjectMap for manager %q", path, managerName)
		}
	}
}

// Tests for ReturnTypeHandler HTTP handler
func TestReturnTypeHandler_MissingParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryParams url.Values
		expectedMsg string
	}{
		{
			name:        "missing name",
			queryParams: url.Values{"type": {"pdf"}, "umbrella": {"false"}},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter",
		},
		{
			name:        "missing type",
			queryParams: url.Values{"name": {"test"}, "umbrella": {"false"}},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter",
		},
		{
			name:        "missing umbrella",
			queryParams: url.Values{"name": {"test"}, "type": {"pdf"}},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter",
		},
		{
			name:        "all missing",
			queryParams: url.Values{},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/?"+test.queryParams.Encode(), nil)
			w := httptest.NewRecorder()

			ReturnTypeHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestReturnTypeHandler_CompositeNotFound(t *testing.T) {
	// Clear Composites
	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{}

	req := httptest.NewRequest("GET", "/?name=nonexistent&type=pdf&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestReturnTypeHandler_ManagerNotInObjectMap(t *testing.T) {
	// Setup folder in Composites but not in ObjectMap
	folder := &Folder{
		Name: "test-folder",
		Files: []*File{
			{Name: "test.pdf", Path: "/test/test.pdf"},
		},
	}

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	// Clear ObjectMap - manager not initialized
	ObjectMap = make(map[string]map[string]object)

	req := httptest.NewRequest("GET", "/?name=test-folder&type=pdf&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array when manager not in ObjectMap
	if len(result) != 0 {
		t.Errorf("Expected 0 files when manager not in ObjectMap, got %d", len(result))
	}
}

func TestReturnTypeHandler_FileTypeFiltering(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()
	managerName := "test-folder-types"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	// Test specific file type filtering
	req := httptest.NewRequest("GET", "/?name=test-folder-types&type=pdf&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return only PDF files
	expectedCount := 1
	if len(result) != expectedCount {
		t.Errorf("Expected %d PDF files, got %d", expectedCount, len(result))
	}

	if len(result) > 0 {
		if result[0].FilePath != "/test/documents/document.pdf" {
			t.Errorf("Expected PDF file path, got %q", result[0].FilePath)
		}
		if result[0].FileName != "document.pdf" {
			t.Errorf("Expected PDF file name, got %q", result[0].FileName)
		}
	}
}

func TestReturnTypeHandler_UmbrellaTypeFiltering(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()
	managerName := "test-folder-types"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	// Test umbrella type filtering for Documents
	req := httptest.NewRequest("GET", "/?name=test-folder-types&type=Documents&umbrella=true", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return Documents (pdf, docx, txt files)
	expectedCount := 3 // document.pdf, report.docx, nested.txt
	if len(result) != expectedCount {
		t.Errorf("Expected %d Document files, got %d", expectedCount, len(result))
	}

	// Verify all returned files are Documents
	for _, file := range result {
		obj := ObjectMap[managerName][file.FilePath]
		if obj.umbrellaType != "Documents" {
			t.Errorf("Expected Documents umbrella type, got %q for file %q", obj.umbrellaType, file.FilePath)
		}
	}
}

func TestReturnTypeHandler_UmbrellaTypeImages(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()
	managerName := "test-folder-types"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	// Test umbrella type filtering for Images
	req := httptest.NewRequest("GET", "/?name=test-folder-types&type=Images&umbrella=true", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return Images (jpg, png, gif files)
	expectedCount := 3 // photo.jpg, logo.png, nested_image.gif
	if len(result) != expectedCount {
		t.Errorf("Expected %d Image files, got %d", expectedCount, len(result))
	}
}

func TestReturnTypeHandler_NonExistentType(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()
	managerName := "test-folder-types"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	// Test filtering for non-existent file type
	req := httptest.NewRequest("GET", "/?name=test-folder-types&type=nonexistent&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array
	if len(result) != 0 {
		t.Errorf("Expected 0 files for non-existent type, got %d", len(result))
	}
}

func TestReturnTypeHandler_EmptyFolder(t *testing.T) {
	// Setup empty folder
	emptyFolder := &Folder{
		Name:  "empty-folder",
		Files: []*File{},
	}

	managerName := "empty-folder"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(emptyFolder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{emptyFolder}

	req := httptest.NewRequest("GET", "/?name=empty-folder&type=pdf&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should return empty array
	if len(result) != 0 {
		t.Errorf("Expected 0 files for empty folder, got %d", len(result))
	}
}

func TestReturnTypeHandler_ResponseFormat(t *testing.T) {
	// Setup test data with one file
	folder := &Folder{
		Name: "response-test",
		Files: []*File{
			{Name: "test.pdf", Path: "/test/test.pdf"},
		},
	}

	managerName := "response-test"

	// Clear and setup ObjectMap and Composites
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	originalComposites := Composites
	defer func() { Composites = originalComposites }()
	Composites = []*Folder{folder}

	req := httptest.NewRequest("GET", "/?name=response-test&type=pdf&umbrella=false", nil)
	w := httptest.NewRecorder()

	ReturnTypeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check Content-Type header
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %q, got %q", expectedContentType, contentType)
	}

	var result []returnStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result))
		return
	}

	// Verify response structure
	file := result[0]
	if file.FilePath != "/test/test.pdf" {
		t.Errorf("Expected file path %q, got %q", "/test/test.pdf", file.FilePath)
	}
	if file.FileName != "test.pdf" {
		t.Errorf("Expected file name %q, got %q", "test.pdf", file.FileName)
	}
}

// Benchmark tests
func BenchmarkGetCategory(b *testing.B) {
	testFiles := []string{
		"document.pdf",
		"image.jpg",
		"music.mp3",
		"video.mp4",
		"unknown.xyz",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range testFiles {
			GetCategory(file)
		}
	}
}

func BenchmarkLoadTypes(b *testing.B) {
	folder := createTestFolderWithTypes()
	managerName := "benchmark-manager"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ObjectMap = make(map[string]map[string]object)
		ObjectMap[managerName] = make(map[string]object)
		LoadTypes(folder, managerName)
	}
}

func BenchmarkReturnTypeHandler(b *testing.B) {
	// Setup
	folder := createTestFolderWithTypes()
	managerName := "benchmark-folder"
	ObjectMap = make(map[string]map[string]object)
	ObjectMap[managerName] = make(map[string]object)
	LoadTypes(folder, managerName)

	// Update folder name to match manager name
	folder.Name = managerName
	Composites = []*Folder{folder}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/?name="+managerName+"&type=Documents&umbrella=true", nil)
		w := httptest.NewRecorder()
		ReturnTypeHandler(w, req)
	}
}
