package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo represents metadata about a stored file
type FileInfo struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	FilePath     string    `json:"file_path"`
}

// StorageConfig holds configuration for file storage
type StorageConfig struct {
	BaseDir      string
	MaxFileSize  int64
	AllowedTypes []string
}

// DefaultStorageConfig returns a default configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		BaseDir:      "./uploads",
		MaxFileSize:  10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{".pdf", ".doc", ".docx", ".txt"},
	}
}

// StorageService handles file operations
type StorageService struct {
	config *StorageConfig
}

// NewStorageService creates a new storage service
func NewStorageService(config *StorageConfig) *StorageService {
	if config == nil {
		config = DefaultStorageConfig()
	}

	// Ensure base directory exists
	if err := os.MkdirAll(config.BaseDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create storage directory: %v", err))
	}

	return &StorageService{config: config}
}

// UploadFile handles file upload with validation and storage
func (s *StorageService) UploadFile(file *multipart.FileHeader, category string) (*FileInfo, error) {
	// Validate file size
	if file.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", file.Size, s.config.MaxFileSize)
	}

	// Validate file type
	if !s.isAllowedFileType(file.Filename) {
		return nil, fmt.Errorf("file type not allowed. Allowed types: %v", s.config.AllowedTypes)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	fileID := s.generateFileID(file.Filename)
	ext := filepath.Ext(file.Filename)
	storedName := fileID + ext

	// Create category directory
	categoryDir := filepath.Join(s.config.BaseDir, category)
	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to create category directory: %v", err)
	}

	// Create the destination file
	filePath := filepath.Join(categoryDir, storedName)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Create file info
	fileInfo := &FileInfo{
		ID:           fileID,
		OriginalName: file.Filename,
		StoredName:   storedName,
		Size:         file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		UploadedAt:   time.Now(),
		FilePath:     filePath,
	}

	return fileInfo, nil
}

// GetFile retrieves a file by ID and category
func (s *StorageService) GetFile(fileID, category string) (*FileInfo, error) {
	categoryDir := filepath.Join(s.config.BaseDir, category)

	// Search for file with the given ID
	files, err := os.ReadDir(categoryDir)
	if err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to read category directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), fileID) {
			filePath := filepath.Join(categoryDir, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Println("error:", err)
				continue
			}

			return &FileInfo{
				ID:           fileID,
				OriginalName: file.Name(),
				StoredName:   file.Name(),
				Size:         fileInfo.Size(),
				UploadedAt:   fileInfo.ModTime(),
				FilePath:     filePath,
			}, nil
		}
	}

	return nil, fmt.Errorf("file not found with ID: %s", fileID)
}

// DeleteFile removes a file by ID and category
func (s *StorageService) DeleteFile(fileID, category string) error {
	fileInfo, err := s.GetFile(fileID, category)
	if err != nil {
		return err
	}

	return os.Remove(fileInfo.FilePath)
}

// ListFiles returns all files in a category
func (s *StorageService) ListFiles(category string) ([]*FileInfo, error) {
	categoryDir := filepath.Join(s.config.BaseDir, category)

	files, err := os.ReadDir(categoryDir)
	if err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to read category directory: %v", err)
	}

	var fileInfos []*FileInfo
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(categoryDir, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Println("error:", err)
				continue
			}

			// Extract file ID from filename (remove extension)
			fileID := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			fileInfos = append(fileInfos, &FileInfo{
				ID:           fileID,
				OriginalName: file.Name(),
				StoredName:   file.Name(),
				Size:         fileInfo.Size(),
				UploadedAt:   fileInfo.ModTime(),
				FilePath:     filePath,
			})
		}
	}

	return fileInfos, nil
}

// isAllowedFileType checks if the file type is allowed
func (s *StorageService) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range s.config.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// generateFileID creates a unique file ID based on filename and timestamp
func (s *StorageService) generateFileID(filename string) string {
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d", filename, timestamp)))
	return hex.EncodeToString(hash[:])
}

// GetFilePath returns the full path to a file
func (s *StorageService) GetFilePath(fileID, category string) (string, error) {
	fileInfo, err := s.GetFile(fileID, category)
	if err != nil {
		return "", err
	}
	return fileInfo.FilePath, nil
}

// FileExists checks if a file exists
func (s *StorageService) FileExists(fileID, category string) bool {
	_, err := s.GetFile(fileID, category)
	return err == nil
}
