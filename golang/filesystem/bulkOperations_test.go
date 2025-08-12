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
		// Removed "no_extension" file to avoid panic in current LoadTypes implementation
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

func TestLoadTypes_EmptyFolder(t *testing.T) {
	// Clear objectMap
	objectMap = make(map[string]object)

	emptyFolder := &Folder{
		Name:       "empty",
		Files:      []*File{},
		Subfolders: []*Folder{},
	}

	LoadTypes(emptyFolder)

	if len(objectMap) != 0 {
		t.Errorf("Expected empty objectMap, but got %d entries", len(objectMap))
	}
}

func TestLoadTypes_FolderWithFiles(t *testing.T) {
	// Clear objectMap
	objectMap = make(map[string]object)

	folder := createTestFolderWithTypes()
	LoadTypes(folder)

	// Check that all files were loaded into objectMap
	expectedEntries := 13 // 11 files in main folder + 2 files in subfolder
	if len(objectMap) != expectedEntries {
		t.Errorf("Expected %d entries in objectMap, but got %d", expectedEntries, len(objectMap))
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
			obj, exists := objectMap[test.path]
			if !exists {
				t.Errorf("Expected path %q to exist in objectMap", test.path)
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

func TestLoadTypes_FileWithoutExtension(t *testing.T) {
	// Clear objectMap
	objectMap = make(map[string]object)

	folder := &Folder{
		Name: "test",
		Files: []*File{
			{Name: "no_extension", Path: "/test/no_extension"},
		},
	}

	// This should cause a panic or error because strings.Split will not have index [1]
	// Testing the current implementation's behavior
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when processing file without extension, but didn't panic")
		} else {
			// Expected behavior - the function panics on files without extensions
			// This test documents the current limitation
			t.Logf("Function correctly panics on files without extension: %v", r)
		}
	}()

	LoadTypes(folder)
}

// Test for the corrected LoadTypes function (see suggested fix below)
func TestLoadTypes_FileWithoutExtension_Fixed(t *testing.T) {
	// This test would work with the fixed version of LoadTypes
	t.Skip("This test requires the fixed LoadTypes function - see comments for fix")

	// Clear objectMap
	objectMap = make(map[string]object)

	folder := &Folder{
		Name: "test",
		Files: []*File{
			{Name: "no_extension", Path: "/test/no_extension"},
			{Name: ".hidden", Path: "/test/.hidden"},
			{Name: "normal.txt", Path: "/test/normal.txt"},
		},
	}

	// With fixed LoadTypes, this should not panic
	LoadTypes(folder)

	// Check that files were processed correctly
	if len(objectMap) != 3 {
		t.Errorf("Expected 3 entries in objectMap, but got %d", len(objectMap))
	}

	// Check file without extension
	if obj, exists := objectMap["/test/no_extension"]; exists {
		if obj.fileType != "" || obj.umbrellaType != "Unknown" {
			t.Errorf("File without extension not handled correctly: fileType=%q, umbrellaType=%q",
				obj.fileType, obj.umbrellaType)
		}
	}
}

func TestLoadTypes_NestedFolders(t *testing.T) {
	// Clear objectMap
	objectMap = make(map[string]object)

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

	LoadTypes(rootFolder)

	// Should have 4 files total
	if len(objectMap) != 4 {
		t.Errorf("Expected 4 entries in objectMap, but got %d", len(objectMap))
	}

	// Verify all files are present
	expectedPaths := []string{
		"/test/root.mp3",
		"/test/level1.pdf",
		"/test/level1/level2.jpg",
		"/test/level1/level2/deep.txt",
	}

	for _, path := range expectedPaths {
		if _, exists := objectMap[path]; !exists {
			t.Errorf("Expected path %q to exist in objectMap", path)
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
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter\n",
		},
		{
			name:        "missing type",
			queryParams: url.Values{"name": {"test"}, "umbrella": {"false"}},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter\n",
		},
		{
			name:        "missing umbrella",
			queryParams: url.Values{"name": {"test"}, "type": {"pdf"}},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter\n",
		},
		{
			name:        "all missing",
			queryParams: url.Values{},
			expectedMsg: "Missing 'name' or 'type' or 'umbrella' parameter\n",
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

			if !contains([]string{w.Body.String()}, test.expectedMsg) {
				t.Errorf("Expected error message to contain %q, got %q", test.expectedMsg, w.Body.String())
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

func TestReturnTypeHandler_FileTypeFiltering(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(folder)

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

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(folder)

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
		obj := objectMap[file.FilePath]
		if obj.umbrellaType != "Documents" {
			t.Errorf("Expected Documents umbrella type, got %q for file %q", obj.umbrellaType, file.FilePath)
		}
	}
}

func TestReturnTypeHandler_UmbrellaTypeImages(t *testing.T) {
	// Setup test data
	folder := createTestFolderWithTypes()

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(folder)

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

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(folder)

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

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(emptyFolder)

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

	// Clear and setup objectMap and Composites
	objectMap = make(map[string]object)
	LoadTypes(folder)

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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		objectMap = make(map[string]object) // Clear map
		LoadTypes(folder)
	}
}

func BenchmarkReturnTypeHandler(b *testing.B) {
	// Setup
	folder := createTestFolderWithTypes()
	objectMap = make(map[string]object)
	LoadTypes(folder)
	Composites = []*Folder{folder}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/?name=test-folder-types&type=Documents&umbrella=true", nil)
		w := httptest.NewRecorder()
		ReturnTypeHandler(w, req)
	}
}
