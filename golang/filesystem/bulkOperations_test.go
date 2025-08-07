package filesystem

import (
	"net/http"
	"net/http/httptest"
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
