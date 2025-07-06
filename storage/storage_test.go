package storage

import (
	"bytes"
	"mime/multipart"
	"os"
	"testing"
)

func TestStorageService_UploadFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create storage service with test config
	config := &StorageConfig{
		BaseDir:      tempDir,
		MaxFileSize:  1024 * 1024, // 1MB
		AllowedTypes: []string{".txt", ".pdf"},
	}

	service := NewStorageService(config)

	// Create a test file content
	testContent := "This is a test file content"

	// Create a mock multipart file
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create a form file
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	// Write content to the form file
	_, err = part.Write([]byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}

	writer.Close()

	// Create a reader for the multipart data
	reader := multipart.NewReader(bytes.NewReader(buf.Bytes()), writer.Boundary())

	// Parse the form
	form, err := reader.ReadForm(1024 * 1024)
	if err != nil {
		t.Fatalf("Failed to read form: %v", err)
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("No file found in form")
	}

	fileHeader := files[0]

	// Test file upload
	fileInfo, err := service.UploadFile(fileHeader, "test")
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}

	// Verify file info
	if fileInfo.OriginalName != "test.txt" {
		t.Errorf("Expected original name 'test.txt', got '%s'", fileInfo.OriginalName)
	}

	if fileInfo.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), fileInfo.Size)
	}

	// Verify file exists
	if !service.FileExists(fileInfo.ID, "test") {
		t.Error("Uploaded file should exist")
	}

	// Verify file content
	filePath, err := service.GetFilePath(fileInfo.ID, "test")
	if err != nil {
		t.Fatalf("Failed to get file path: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read uploaded file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}
}

func TestStorageService_FileTypeValidation(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{
		BaseDir:      tempDir,
		MaxFileSize:  1024 * 1024,
		AllowedTypes: []string{".txt", ".pdf"},
	}

	service := NewStorageService(config)

	// Test allowed file type
	if !service.isAllowedFileType("document.txt") {
		t.Error("'.txt' should be an allowed file type")
	}

	// Test disallowed file type
	if service.isAllowedFileType("image.jpg") {
		t.Error("'.jpg' should not be an allowed file type")
	}
}

func TestStorageService_FileSizeValidation(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{
		BaseDir:      tempDir,
		MaxFileSize:  100, // 100 bytes
		AllowedTypes: []string{".txt"},
	}

	service := NewStorageService(config)

	// Create a large content that exceeds the size limit
	largeContent := make([]byte, 200) // 200 bytes

	// Create a mock multipart file
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "large.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = part.Write(largeContent)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}

	writer.Close()

	reader := multipart.NewReader(bytes.NewReader(buf.Bytes()), writer.Boundary())
	form, err := reader.ReadForm(1024 * 1024)
	if err != nil {
		t.Fatalf("Failed to read form: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("No file found in form")
	}

	fileHeader := files[0]

	// Test that upload fails due to size limit
	_, err = service.UploadFile(fileHeader, "test")
	if err == nil {
		t.Error("Should fail when file size exceeds limit")
	}
}
