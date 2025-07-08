package api

import (
	"cold_emailer/constants"
	"cold_emailer/dbmodels/profileinfo"
	"cold_emailer/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cold_emailer/gmail"
	"cold_emailer/openai"

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
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_FORM",
			Message: "Failed to parse form data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var profileReq ProfileUploadRequest
	if err := c.ShouldBind(&profileReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_PROFILE_DATA",
			Message: fmt.Sprintf("Invalid profile data: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Handle resume file upload
	var resumeFile *storage.FileInfo
	var resumePath string
	var resumeMetadata string
	if file, err := c.FormFile("resume"); err == nil && file != nil {
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
		resumePath = uploadedFile.FilePath
		// Save resume metadata as JSON using a map
		meta := map[string]interface{}{
			"resume_metadata": map[string]interface{}{
				"original_name": uploadedFile.OriginalName,
				"stored_name":   uploadedFile.StoredName,
				"size":          uploadedFile.Size,
				"mime_type":     uploadedFile.MimeType,
				"uploaded_at":   uploadedFile.UploadedAt.Format(time.RFC3339),
			},
		}
		b, err := json.Marshal(meta)
		if err == nil {
			resumeMetadata = string(b)
		} else {
			resumeMetadata = "{}"
		}
	}

	profileID := uuid.New().String()

	// Insert into DB
	info := profileinfo.StructForSet{
		ID:          profileID,
		Status:      "ACTIVE",
		Name:        profileReq.Name,
		Email:       profileReq.Email,
		Phone:       profileReq.Phone,
		LinkedInURL: profileReq.LinkedInURL,
		Experience:  profileReq.Experience,
		Skills:      profileReq.Skills,
		Summary:     profileReq.Summary,
		ResumePath:  resumePath,
		Metadata:    resumeMetadata,
	}
	if err := profileinfo.Create(context.Background(), info); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DB_ERROR",
			Message: fmt.Sprintf("Failed to save profile: %v", err),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := ProfileUploadResponse{
		Message:   "Profile uploaded successfully",
		ProfileID: profileID,
		CreatedAt: time.Now(),
	}
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

	// Fetch latest active profile from DB
	profile, err := profileinfo.GetLatestActive(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "NO_PROFILE",
			Message: "No active profile found. Please upload your profile first.",
			Code:    http.StatusBadRequest,
		})
		return
	}

	profileText := fmt.Sprintf(
		"Name: %s\nEmail: %s\nPhone: %s\nLinkedIn: %s\nExperience: %s\nSkills: %s\nSummary: %s",
		profile.Name,
		profile.Email,
		profile.Phone,
		profile.LinkedInURL,
		profile.Experience,
		profile.Skills,
		profile.Summary,
	)

	prompt := fmt.Sprintf(
		"Generate a personalized cold outreach email for a job opportunity. Company Name: %s\n\nProfile:\n%s\n\nTarget Info: %s\n\n%s",
		req.CompanyName,
		profileText,
		req.TargetID, // In a real app, you'd look up the target info by ID
		req.CustomPrompt,
	)

	openaiClient := openai.NewOpenAIClient()
	emailText, err := openaiClient.GenerateEmail(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "OPENAI_ERROR",
			Message: fmt.Sprintf("Failed to generate email: %v", err),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := GenerateEmailResponse{
		Message:     "Email generated successfully",
		EmailDraft:  emailText,
		Subject:     "Personalized Cold Outreach Email",
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

// Gmail OAuth2: Initiate auth flow
func GmailAuthInitiateHandler(c *gin.Context) {
	state := "state-token" // In production, generate a random state and validate it in callback
	url := gmail.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"auth_url": url})
}

// Gmail OAuth2: Callback
func GmailOAuth2CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code in callback"})
		return
	}
	ctx := context.Background()
	tok, err := gmail.ExchangeCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// For now, return the tokens in the response (in production, store securely)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  tok.AccessToken,
		"refresh_token": tok.RefreshToken,
		"expiry":        tok.Expiry,
	})
}
