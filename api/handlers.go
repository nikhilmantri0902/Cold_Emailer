package api

import (
	"cold_emailer/constants"
	"cold_emailer/storage"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Global storage service instance
var storageService *storage.StorageService

// InitStorage initializes the storage service
func InitStorage() {
	storageService = storage.NewStorageService(nil) // Use default config
}

// Upload user profile and resume
func UploadProfileHandler(c *gin.Context) {
	// Parse multipart form (max 32MB)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_FORM",
			Message: "Failed to parse form data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Parse profile data from form
	var profileReq ProfileUploadRequest
	if err := c.ShouldBind(&profileReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_PROFILE_DATA",
			Message: fmt.Sprintf("Invalid profile data: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// TODO: Store the profile data in the database
	log.Println("profileReq", profileReq)

	// Handle resume file upload
	var resumeFile *storage.FileInfo
	if file, err := c.FormFile("resume"); err == nil && file != nil {
		// Upload resume file
		uploadedFile, err := storageService.UploadFile(file, constants.RESUME_CATEGORY)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "RESUME_UPLOAD_FAILED",
				Message: fmt.Sprintf("Failed to upload resume: %v", err),
				Code:    http.StatusBadRequest,
			})
			return
		}
		resumeFile = uploadedFile
	}

	// Generate profile ID
	profileID := uuid.New().String()

	// Create response
	response := ProfileUploadResponse{
		Message:   "Profile uploaded successfully",
		ProfileID: profileID,
		CreatedAt: time.Now(),
	}

	// Add resume file info if uploaded
	if resumeFile != nil {
		response.ResumeFile = &FileInfo{
			ID:           resumeFile.ID,
			OriginalName: resumeFile.OriginalName,
			StoredName:   resumeFile.StoredName,
			Size:         resumeFile.Size,
			MimeType:     resumeFile.MimeType,
			UploadedAt:   resumeFile.UploadedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// Upload targets (CSV or JSON list)
func UploadTargetsHandler(c *gin.Context) {
	var req TargetsUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_TARGETS_DATA",
			Message: fmt.Sprintf("Invalid targets data: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Generate target IDs
	var targetIDs []string
	for range req.Targets {
		targetIDs = append(targetIDs, uuid.New().String())
	}

	response := TargetsUploadResponse{
		Message:   fmt.Sprintf("Successfully uploaded %d targets", len(req.Targets)),
		TargetIDs: targetIDs,
		Count:     len(req.Targets),
	}

	c.JSON(http.StatusOK, response)
}

// Generate personalized email for a target
func GenerateEmailHandler(c *gin.Context) {
	var req GenerateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: fmt.Sprintf("Invalid request data: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// TODO: Call OpenAI, generate and return email
	response := GenerateEmailResponse{
		Message:     "Email generated successfully",
		EmailDraft:  "This is a placeholder email draft. OpenAI integration pending.",
		Subject:     "Interested in discussing opportunities",
		GeneratedAt: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// Send email to a target (attach resume)
func SendEmailHandler(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: fmt.Sprintf("Invalid request data: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// TODO: Implement Gmail send logic
	response := SendEmailResponse{
		Message: "Email sent successfully",
		EmailID: uuid.New().String(),
		SentAt:  time.Now(),
		Status:  "sent",
	}

	c.JSON(http.StatusOK, response)
}

// Get status of sent emails
func StatusHandler(c *gin.Context) {
	// TODO: Return email send status from DB
	response := StatusResponse{
		Status: "ok",
		Emails: []EmailStatus{
			{
				EmailID:     "sample-email-id",
				TargetID:    "sample-target-id",
				TargetName:  "John Doe",
				TargetEmail: "john@example.com",
				Company:     "Example Corp",
				Status:      "sent",
				SentAt:      time.Now(),
			},
		},
		Total:   1,
		Success: 1,
		Failed:  0,
	}

	c.JSON(http.StatusOK, response)
}

// Get logs (basic)
func LogsHandler(c *gin.Context) {
	// TODO: Return recent operation logs
	response := LogsResponse{
		Logs: []LogEntry{
			{
				ID:        "log-1",
				Level:     "info",
				Message:   "Profile uploaded successfully",
				Timestamp: time.Now(),
				Category:  "upload",
			},
		},
		Total: 1,
	}

	c.JSON(http.StatusOK, response)
}
