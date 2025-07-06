package api

import "time"

// Profile upload request/response
type ProfileUploadRequest struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Email       string `json:"email" form:"email" binding:"required,email"`
	Phone       string `json:"phone,omitempty" form:"phone,omitempty"`
	LinkedInURL string `json:"linkedin_url,omitempty" form:"linkedin_url,omitempty"`
	Experience  string `json:"experience,omitempty" form:"experience,omitempty"`
	Skills      string `json:"skills,omitempty" form:"skills,omitempty"`
	Summary     string `json:"summary,omitempty" form:"summary,omitempty"`
	// Resume file will be handled separately via multipart form
}

type ProfileUploadResponse struct {
	Message    string    `json:"message"`
	ProfileID  string    `json:"profile_id"`
	ResumeFile *FileInfo `json:"resume_file,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// FileInfo represents uploaded file metadata
type FileInfo struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// Target upload request/response
type Target struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Company     string `json:"company" binding:"required"`
	Role        string `json:"role" binding:"required"` // e.g. CTO, Recruiter
	Description string `json:"description,omitempty"`   // Optional
	LinkedInURL string `json:"linkedin_url,omitempty"`
}

type TargetsUploadRequest struct {
	Targets []Target `json:"targets" binding:"required,min=1"`
}

type TargetsUploadResponse struct {
	Message   string   `json:"message"`
	TargetIDs []string `json:"target_ids,omitempty"`
	Count     int      `json:"count"`
}

// Email generation request/response
type GenerateEmailRequest struct {
	TargetID     string `json:"target_id" binding:"required"`
	Template     string `json:"template,omitempty"`      // Optional template name
	CustomPrompt string `json:"custom_prompt,omitempty"` // Optional custom prompt
}

type GenerateEmailResponse struct {
	Message     string    `json:"message"`
	EmailDraft  string    `json:"email_draft,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
}

// Email send request/response
type SendEmailRequest struct {
	TargetID      string `json:"target_id" binding:"required"`
	EmailBody     string `json:"email_body" binding:"required"`
	Subject       string `json:"subject" binding:"required"`
	IncludeResume bool   `json:"include_resume,omitempty"` // Whether to attach resume
}

type SendEmailResponse struct {
	Message string    `json:"message"`
	EmailID string    `json:"email_id,omitempty"`
	SentAt  time.Time `json:"sent_at"`
	Status  string    `json:"status"` // "sent", "failed", "pending"
}

// Status/logs response
type EmailStatus struct {
	EmailID     string    `json:"email_id"`
	TargetID    string    `json:"target_id"`
	TargetName  string    `json:"target_name"`
	TargetEmail string    `json:"target_email"`
	Company     string    `json:"company"`
	Status      string    `json:"status"` // "sent", "failed", "pending"
	SentAt      time.Time `json:"sent_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

type StatusResponse struct {
	Status  string        `json:"status"`
	Emails  []EmailStatus `json:"emails"`
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
}

type LogEntry struct {
	ID        string    `json:"id"`
	Level     string    `json:"level"` // "info", "error", "warning"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Category  string    `json:"category"` // "upload", "email", "generation"
}

type LogsResponse struct {
	Logs  []LogEntry `json:"logs"`
	Total int        `json:"total"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
