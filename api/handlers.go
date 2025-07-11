package api

import (
	"cold_emailer/constants"
	"cold_emailer/dbmodels/gmailtokens"
	"cold_emailer/dbmodels/profileinfo"
	"cold_emailer/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cold_emailer/dbmodels/contacts"
	"cold_emailer/dbmodels/email_logs"
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
		log.Println("error:", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_FORM",
			Message: "Failed to parse form data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var profileReq ProfileUploadRequest
	if err := c.ShouldBind(&profileReq); err != nil {
		log.Println("error:", err)
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
			log.Println("error:", err)
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
			log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error:", err)
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
		log.Println("error: missing code in callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code in callback"})
		return
	}
	ctx := context.Background()
	tok, err := gmail.ExchangeCode(ctx, code)
	if err != nil {
		log.Println("error:", err)
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
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Gmail token found. Please authenticate first."})
		return
	}
	from := token.EmailID
	if from == "me" {
		from = "your@email.com" // fallback, replace with your real email if needed
	}
	err = gmail.SendSingleEmail(ctx, token.AccessToken, from, req.To, req.Subject, req.Body)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// EnrichDatabaseHandler triggers Apollo enrichment and returns the result
func EnrichDatabaseHandler(c *gin.Context) {
	go func() {
		err := EnrichDBWithCompaniesAndContacts(context.Background(), 10, 100)
		if err != nil {
			log.Println("err:", err)
			return
		}

		log.Println("Enrichment complete")
	}()
	c.JSON(http.StatusAccepted, gin.H{"message": "Enrichment started. Check logs for progress."})
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

	// Start email sending in goroutine
	go func() {
		ctx := context.Background()
		log.Printf("Starting to send %d initial emails to contacts with status: %s", req.Count, req.Status)

		// Step 2 - Fetch contacts with company info
		contactsList, err := contacts.GetContactsWithCompanyInfo(ctx, req.Count, req.Status, "")
		if err != nil {
			log.Printf("ERROR: Failed to fetch contacts: %v", err)
			return
		}

		log.Printf("Found %d contacts to send emails to", len(contactsList))

		// Get profile info for email generation
		profile, err := profileinfo.GetLatestActive(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to get profile info: %v", err)
			return
		}

		// Get Gmail token
		token, err := gmailtokens.GetLatestToken(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to get Gmail token: %v", err)
			return
		}

		from := token.EmailID
		if from == "me" {
			from = profile.Email // Use profile email as fallback
		}

		// Initialize OpenAI client
		openaiClient := openai.NewOpenAIClient()

		// Track results
		emailsGenerated := 0
		emailsSent := 0
		emailsFailed := 0

		// Process each contact
		for _, contact := range contactsList {
			log.Printf("Processing contact: %s (%s) at %s", contact.ContactName, contact.ContactEmail, contact.CompanyName)

			// Step 3: Generate personalized email
			emailData := openai.EmailGenerationData{
				ContactName:       contact.ContactName,
				ContactRole:       contact.ContactRole,
				ContactLinkedIn:   contact.ContactLinkedIn,
				CompanyName:       contact.CompanyName,
				CompanyWebsite:    contact.CompanyWebsite,
				CompanyIndustry:   contact.CompanyIndustry,
				ProfileName:       profile.Name,
				ProfileExperience: profile.Experience,
				ProfileSkills:     profile.Skills,
				ProfileSummary:    profile.Summary,
			}

			subject, body, err := openaiClient.GeneratePersonalizedEmail(emailData)
			if err != nil {
				log.Printf("ERROR: Failed to generate email for %s: %v", contact.ContactEmail, err)
				emailsFailed++
				continue
			}

			emailsGenerated++
			log.Printf("Generated email for %s: %s", contact.ContactEmail, subject)

			// Log GENERATED stage
			metadata := map[string]interface{}{
				"contact_name": contact.ContactName,
				"contact_role": contact.ContactRole,
				"company_name": contact.CompanyName,
				"generated_at": time.Now().Format(time.RFC3339),
			}

			if err := email_logs.LogGenerated(ctx, contact.ContactID, contact.CompanyID, subject, body, metadata); err != nil {
				log.Printf("ERROR: Failed to log GENERATED stage for %s: %v", contact.ContactEmail, err)
			}

			// Step 4: Send email with resume
			if profile.ResumePath == "" {
				log.Printf("WARNING: No resume path found for profile, sending without attachment")
				err = gmail.SendSingleEmail(ctx, token.AccessToken, from, contact.ContactEmail, subject, body)
			} else {
				err = gmail.SendEmailWithAttachment(ctx, token.AccessToken, from, contact.ContactEmail, subject, body, profile.ResumePath)
			}

			if err != nil {
				log.Printf("ERROR: Failed to send email to %s: %v", contact.ContactEmail, err)
				emailsFailed++

				// Log ERROR stage
				errorMetadata := map[string]interface{}{
					"contact_name": contact.ContactName,
					"contact_role": contact.ContactRole,
					"company_name": contact.CompanyName,
					"error_at":     time.Now().Format(time.RFC3339),
					"error_type":   "send_failed",
				}

				if logErr := email_logs.LogError(ctx, contact.ContactID, contact.CompanyID, subject, body, err.Error(), errorMetadata); logErr != nil {
					log.Printf("ERROR: Failed to log ERROR stage for %s: %v", contact.ContactEmail, logErr)
				}
				continue
			}

			emailsSent++
			log.Printf("SUCCESS: Sent email to %s (%s)", contact.ContactEmail, subject)

			// Log SENT stage
			sentMetadata := map[string]interface{}{
				"contact_name":    contact.ContactName,
				"contact_role":    contact.ContactRole,
				"company_name":    contact.CompanyName,
				"sent_at":         time.Now().Format(time.RFC3339),
				"resume_attached": profile.ResumePath != "",
			}

			if err := email_logs.LogSent(ctx, contact.ContactID, contact.CompanyID, subject, body, sentMetadata); err != nil {
				log.Printf("ERROR: Failed to log SENT stage for %s: %v", contact.ContactEmail, err)
			}

			// Small delay to avoid rate limiting
			time.Sleep(2 * time.Second)
		}

		log.Printf("Email sending completed: %d generated, %d sent, %d failed", emailsGenerated, emailsSent, emailsFailed)
	}()

	c.JSON(http.StatusAccepted, SendFewInitialEmailsResponse{
		Message: "Email sending started. Check logs for progress.",
		Count:   req.Count,
	})
}
