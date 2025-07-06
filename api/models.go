package api

// Profile upload request/response
type ProfileUploadResponse struct {
	Message string `json:"message"`
	FileURL string `json:"file_url,omitempty"`
}

// Target upload request/response
type Target struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`        // e.g. CTO, Recruiter
	Description string `json:"description"` // Optional
}

type TargetsUploadRequest struct {
	Targets []Target `json:"targets"`
}

type TargetsUploadResponse struct {
	Message   string   `json:"message"`
	TargetIDs []string `json:"target_ids,omitempty"`
}

// Email generation request/response
type GenerateEmailRequest struct {
	TargetID string `json:"target_id"`
	// Add more fields if needed
}

type GenerateEmailResponse struct {
	Message    string `json:"message"`
	EmailDraft string `json:"email_draft,omitempty"`
}

// Email send request/response
type SendEmailRequest struct {
	TargetID  string `json:"target_id"`
	EmailBody string `json:"email_body"`
}

type SendEmailResponse struct {
	Message string `json:"message"`
}

// Status/logs response
type StatusResponse struct {
	Status  string   `json:"status"`
	Details []string `json:"details"`
}

type LogsResponse struct {
	Logs []string `json:"logs"`
}
