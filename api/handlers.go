package api

import (
	"cold_emailer/constants"
	"cold_emailer/dbmodels/companies"
	"cold_emailer/dbmodels/contacts"
	"cold_emailer/dbmodels/email_logs"
	"cold_emailer/dbmodels/gmailtokens"
	"cold_emailer/dbmodels/jobs"
	"cold_emailer/dbmodels/profileinfo"
	"cold_emailer/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cold_emailer/gmail"
	"cold_emailer/openai"

	"cold_emailer/utils"

	"database/sql"

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

		ctx := context.Background()
		jobID := utils.GenerateUUID()
		var err error
		jobID, err = jobs.Insert(ctx, jobs.JobForSet{
			ID:     jobID,
			Name:   constants.JobNameEnrichCompaniesAndContacts,
			Status: jobs.StatusActive,
			Metadata: utils.MarshalInterface(
				map[string]interface{}{
					"request_params": map[string]interface{}{
						"company_count":            companyCount,
						"max_contacts_per_company": maxContactsPerCompany,
					},
				},
			),
		})

		defer func(jobID string) {
			status := jobs.StatusCompleted
			message := "Enrichment complete"

			if err != nil {
				status = jobs.StatusFailed
				message = fmt.Sprintf("Enrichment failed: %v", err)
			}

			jobs.Update(ctx, jobs.UpdateInput{
				ID:      jobID,
				Status:  &status,
				Message: &message,
			})
		}(jobID)

		err = EnrichDBWithCompaniesAndContacts(ctx, companyCount, maxContactsPerCompany)
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
		ctx := context.Background()
		jobID := utils.GenerateUUID()
		var err error
		jobID, err = jobs.Insert(ctx, jobs.JobForSet{
			ID:     jobID,
			Name:   constants.JobNameBackfillCompanyDetails,
			Status: jobs.StatusActive,
			Metadata: utils.MarshalInterface(
				map[string]interface{}{
					"request_params": map[string]interface{}{
						"operation": "backfill_company_details",
					},
				},
			),
		})

		defer func(jobID string) {
			status := jobs.StatusCompleted
			message := "Backfill complete"

			if err != nil {
				status = jobs.StatusFailed
				message = fmt.Sprintf("Backfill failed: %v", err)
			}

			jobs.Update(ctx, jobs.UpdateInput{
				ID:      jobID,
				Status:  &status,
				Message: &message,
			})
		}(jobID)

		err = BackFillCompanyDetailsFunc(ctx)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
		utils.Logger.Info().Msg("Backfill complete")
	}()
	c.JSON(http.StatusAccepted, gin.H{"message": "Backfilling started. Check logs for progress."})
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
		jobID := utils.GenerateUUID()
		var err error
		jobID, err = jobs.Insert(ctx, jobs.JobForSet{
			ID:     jobID,
			Name:   constants.JobNameSendEmails,
			Status: jobs.StatusActive,
			Metadata: utils.MarshalInterface(
				map[string]interface{}{
					"request_params": map[string]interface{}{
						"count":  req.Count,
						"status": req.Status,
					},
				},
			),
		})

		defer func(jobID string) {
			status := jobs.StatusCompleted
			message := "Email sending complete"

			if err != nil {
				status = jobs.StatusFailed
				message = fmt.Sprintf("Email sending failed: %v", err)
			}

			jobs.Update(ctx, jobs.UpdateInput{
				ID:      jobID,
				Status:  &status,
				Message: &message,
			})
		}(jobID)

		utils.Logger.Info().Int("count", req.Count).Str("status", req.Status).Msg("Starting to send initial emails to contacts with status:")

		err = GenerateAndSendEmails(ctx, req.Count, req.Status)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
		utils.Logger.Info().Msg("Email sending complete")
	}()

	c.JSON(http.StatusAccepted, SendFewInitialEmailsResponse{
		Message: "Email sending started. Check logs for progress.",
		Count:   req.Count,
	})
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
		jobID := utils.GenerateUUID()
		var err error
		jobID, err = jobs.Insert(ctx, jobs.JobForSet{
			ID:     jobID,
			Name:   constants.JobNameSendEmailFollowUp,
			Status: jobs.StatusActive,
			Metadata: utils.MarshalInterface(
				map[string]interface{}{
					"request_params": map[string]interface{}{
						"count":                 req.Count,
						"status":                req.Status,
						"days_past_first_email": req.DaysPastFirstEmail,
					},
				},
			),
		})

		defer func(jobID string) {
			status := jobs.StatusCompleted
			message := "Follow up email sending complete"

			if err != nil {
				status = jobs.StatusFailed
				message = fmt.Sprintf("Follow up email sending failed: %v", err)
			}

			jobs.Update(ctx, jobs.UpdateInput{
				ID:      jobID,
				Status:  &status,
				Message: &message,
			})
		}(jobID)

		utils.Logger.Info().Int("count", req.Count).Str("status", req.Status).Msg("Starting to send follow up emails to contacts with status:")

		err = GenerateAndSendFollowUpEmails(ctx, req.DaysPastFirstEmail, req.Count, req.Status)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return
		}
		utils.Logger.Info().Msg("Follow up email sending complete")
	}()

	c.JSON(http.StatusAccepted, SendFewFollowUpEmailsResponse{
		Message: "Follow up email sending started. Check logs for progress.",
		Count:   req.Count,
	})
}

// GetEmailLogsHandler returns paginated email logs with company and contact info
func GetEmailLogsHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")
	stage := c.Query("stage") // Optional filter by email stage

	pageInt := 1
	pageSizeInt := 50

	if p, err := strconv.Atoi(page); err == nil {
		pageInt = p
	}
	if ps, err := strconv.Atoi(pageSize); err == nil {
		pageSizeInt = ps
	}

	offset := (pageInt - 1) * pageSizeInt

	ctx := context.Background()
	logs, total, err := email_logs.GetEmailLogsWithPagination(ctx, offset, pageSizeInt, stage)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"pagination": gin.H{
			"page":        pageInt,
			"page_size":   pageSizeInt,
			"total":       total,
			"total_pages": (total + pageSizeInt - 1) / pageSizeInt,
		},
	})
}

// GetCompaniesHandler returns all companies
func GetCompaniesHandler(c *gin.Context) {
	ctx := context.Background()
	companies, err := companies.GetAll(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"total":     len(companies),
	})
}

// GetCompanyContactsHandler returns contacts for a specific company
func GetCompanyContactsHandler(c *gin.Context) {
	companyID := c.Param("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	ctx := context.Background()
	contactsList, err := contacts.GetContactsByCompanyID(ctx, companyID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contacts": contactsList,
		"total":    len(contactsList),
	})
}

// GetSystemConfigHandler returns system configuration including target countries and filters
func GetSystemConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"target_countries": constants.TargetCountries,
		"suitable_roles":   constants.SuitableRoles,
		"company_sizes":    constants.OrganizationEmployeeRanges,
		"industries":       constants.OrganizationKeywordTags,
	})
}

// GetJobStatusHandler returns the status of a job by name
func GetJobStatusHandler(c *gin.Context) {
	jobName := c.Query("job_name")
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_name query parameter is required"})
		return
	}

	ctx := context.Background()
	job, err := jobs.GetByName(ctx, jobName, "")
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get job status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_id":       job.ID,
		"job_name":     job.Name,
		"status":       job.Status,
		"message":      job.Message,
		"created_at":   job.CreatedAt,
		"completed_at": job.CompletedAt,
		"metadata":     job.Metadata,
	})
}
