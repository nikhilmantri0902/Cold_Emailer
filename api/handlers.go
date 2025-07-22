package api

import (
	"cold_emailer/constants"
	"cold_emailer/dbmodels/gmailtokens"
	"cold_emailer/dbmodels/profileinfo"
	"cold_emailer/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cold_emailer/dbmodels/contacts"
	"cold_emailer/gmail"
	"cold_emailer/openai"

	"cold_emailer/utils"

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
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_FORM",
			Message: "Failed to parse form data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var profileReq ProfileUploadRequest
	if err := c.ShouldBind(&profileReq); err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
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
			utils.Logger.Error().Err(err).Msg("error:")
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
			utils.Logger.Error().Err(err).Msg("error:")
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
		utils.Logger.Error().Err(err).Msg("error:")
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

// Generate personalized email for a target
func GenerateEmailHandler(c *gin.Context) {
	var req GenerateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
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
		utils.Logger.Error().Err(err).Msg("error:")
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
		utils.Logger.Error().Err(err).Msg("error:")
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
		utils.Logger.Error().Msg("error: missing code in callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code in callback"})
		return
	}
	ctx := context.Background()
	tok, err := gmail.ExchangeCode(ctx, code)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For now, get the user's email from the token info (in production, use Gmail API to get profile)
	emailID := c.Query("email_id") // Optionally pass email_id as query param
	if emailID == "" {
		emailID = "me" // fallback, but should be set properly
	}

	tokenRow := gmailtokens.GmailTokenForSet{
		ID:           uuid.New().String(),
		EmailID:      emailID,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry.Format(time.RFC3339),
	}
	_ = gmailtokens.CreateGmailToken(ctx, tokenRow)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tok.AccessToken,
		"refresh_token": tok.RefreshToken,
		"expiry":        tok.Expiry,
	})
}

// SendSingleEmailRequest is the request body for /send-single-email
type SendSingleEmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// SendSingleEmailHandler handles sending a single test email
func SendSingleEmailHandler(c *gin.Context) {
	var req SendSingleEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()

	token, err := gmailtokens.GetLatestToken(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Gmail token found. Please authenticate first."})
		return
	}

	// check if token has expired, return reauthenticate error
	if token.CheckExpiry() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gmail oauth credentials have expired. Please reauthenticate"})
		return
	}

	from := token.EmailID
	if from == "me" {
		from = "your@email.com" // fallback, replace with your real email if needed
	}
	err = gmail.SendSingleEmail(ctx, token.AccessToken, from, req.To, req.Subject, req.Body)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// EnrichDatabaseHandler triggers Apollo enrichment and returns the result
func EnrichDatabaseHandler(c *gin.Context) {
	var req struct {
		CompanyCount          int `json:"company_count" binding:"min=1,max=100"`
		MaxContactsPerCompany int `json:"max_contacts_per_company" binding:"min=1,max=1000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyCount := req.CompanyCount
	if companyCount == 0 {
		companyCount = 10
	}
	maxContactsPerCompany := req.MaxContactsPerCompany
	if maxContactsPerCompany == 0 {
		maxContactsPerCompany = 100
	}

	utils.Logger.Info().Int("company_count", companyCount).Int("max_contacts_per_company", maxContactsPerCompany).Msg("company_count:")
	go func() {
		err := EnrichDBWithCompaniesAndContacts(context.Background(), companyCount, maxContactsPerCompany)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
		utils.Logger.Info().Msg("Enrichment complete")
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Enrichment started. Check logs for progress."})
}

// EnrichDatabaseHandler triggers Apollo enrichment and returns the result
func BackfillCompanyDetails(c *gin.Context) {
	go func() {
		err := BackFillCompanyDetailsFunc(context.Background())
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{"message": "Backfilling started. Check logs for progress."})
}

// SendFewInitialEmailsRequest is the request body for sending initial emails
type SendFewInitialEmailsRequest struct {
	Count  int    `json:"count" binding:"min=1,max=50"`
	Status string `json:"status,omitempty"` // Optional status filter for contacts
}

// SendFewInitialEmailsResponse is the response for the initial emails endpoint
type SendFewInitialEmailsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// SendFewInitialEmailsHandler sends initial emails to contacts
func SendFewInitialEmailsHandler(c *gin.Context) {
	var req SendFewInitialEmailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default count if not provided
	if req.Count == 0 {
		req.Count = 10
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = contacts.StatusActive
	}

	// Get Gmail token
	token, err := gmailtokens.GetLatestToken(c.Request.Context())
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Gmail token found. Please authenticate first."})
		return
	}

	// check if token has expired, return reauthenticate error
	if token.CheckExpiry() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gmail oauth credentials have expired. Please reauthenticate"})
		return
	}

	// Start email sending in goroutine
	go func() {
		ctx := context.Background()
		utils.Logger.Info().Int("count", req.Count).Str("status", req.Status).Msg("Starting to send initial emails to contacts with status:")

		err = GenerateAndSendEmails(ctx, req.Count, req.Status)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
	}()

	c.JSON(http.StatusAccepted, SendFewInitialEmailsResponse{
		Message: "Email sending started. Check logs for progress.",
		Count:   req.Count,
	})
}

// SendFewFollowUpmailsRequest is the request body for sending follow up emails
type SendFewFollowUpmailsRequest struct {
	Count              int    `json:"count" binding:"min=1,max=50"`
	Status             string `json:"status,omitempty"` // Optional status filter for contacts
	DaysPastFirstEmail int    `json:"days_past_first_email,omitempty"`
}

// SendFewFollowUpEmailsResponse is the response for the follow up emails endpoint
type SendFewFollowUpEmailsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// SendFewFollowUpEmailsHandler sends follow_up emails to contacts
func SendFewFollowUpEmailsHandler(c *gin.Context) {
	var req SendFewFollowUpmailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default count if not provided
	if req.Count == 0 {
		req.Count = 10
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = contacts.StatusActive
	}

	if req.DaysPastFirstEmail == 0 {
		req.DaysPastFirstEmail = 3
	}

	// Get Gmail token
	token, err := gmailtokens.GetLatestToken(c.Request.Context())
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Gmail token found. Please authenticate first."})
		return
	}

	// check if token has expired, return reauthenticate error
	if token.CheckExpiry() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gmail oauth credentials have expired. Please reauthenticate"})
		return
	}

	// Start email sending in goroutine
	go func() {
		ctx := context.Background()
		utils.Logger.Info().Int("count", req.Count).Str("status", req.Status).Msg("Starting to send follow up emails to contacts with status:")

		err = GenerateAndSendFollowUpEmails(ctx, req.DaysPastFirstEmail, req.Count, req.Status)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
	}()

	c.JSON(http.StatusAccepted, SendFewFollowUpEmailsResponse{
		Message: "Follow up email sending started. Check logs for progress.",
		Count:   req.Count,
	})
}
