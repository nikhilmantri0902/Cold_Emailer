package email_logs

import (
	"cold_emailer/db"
	"cold_emailer/utils"
	"context"
	"encoding/json"
)

// InsertLog inserts a new email log entry
func InsertLog(ctx context.Context, l EmailLogForSet) error {
	query := `INSERT INTO email_logs (id, contact_id, company_id, status, email_stage, email_subject, email_body, attachment_details, error_message, metadata, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())`
	_, err := db.GetDB().ExecContext(ctx, query, l.ID, l.ContactID, l.CompanyID, l.Status, l.EmailStage, l.EmailSubject, l.EmailBody, l.AttachmentDetails, l.ErrorMessage, l.Metadata)
	return err
}

// LogGenerated logs when an email is generated
func LogGenerated(ctx context.Context, contactID, companyID, subject, body string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	log := EmailLogForSet{
		ID:                utils.GenerateUUID(),
		ContactID:         contactID,
		CompanyID:         companyID,
		Status:            StatusActive,
		EmailStage:        "GENERATED",
		EmailSubject:      subject,
		EmailBody:         body,
		AttachmentDetails: "{}",
		Metadata:          string(metadataJSON),
	}
	return InsertLog(ctx, log)
}

// LogGenerated logs when an email is generated
func LogGeneratedFollowUp(ctx context.Context, contactID, companyID, subject, body string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	log := EmailLogForSet{
		ID:                utils.GenerateUUID(),
		ContactID:         contactID,
		CompanyID:         companyID,
		Status:            StatusActive,
		EmailStage:        "GENERATED_FOLLOW_UP",
		EmailSubject:      subject,
		EmailBody:         body,
		AttachmentDetails: "{}",
		Metadata:          string(metadataJSON),
	}
	return InsertLog(ctx, log)
}

// LogSent logs when an email is successfully sent
func LogSent(ctx context.Context, contactID, companyID, subject, body string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	log := EmailLogForSet{
		ID:                utils.GenerateUUID(),
		ContactID:         contactID,
		CompanyID:         companyID,
		Status:            StatusActive,
		EmailStage:        "SENT",
		EmailSubject:      subject,
		EmailBody:         body,
		AttachmentDetails: "{}",
		Metadata:          string(metadataJSON),
	}
	return InsertLog(ctx, log)
}

// LogSent logs when an email is successfully sent
func LogSentFollowUp(ctx context.Context, contactID, companyID, subject, body string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	log := EmailLogForSet{
		ID:                utils.GenerateUUID(),
		ContactID:         contactID,
		CompanyID:         companyID,
		Status:            StatusActive,
		EmailStage:        "SENT_FOLLOW_UP",
		EmailSubject:      subject,
		EmailBody:         body,
		AttachmentDetails: "{}",
		Metadata:          string(metadataJSON),
	}
	return InsertLog(ctx, log)
}

// LogError logs when an email fails
func LogError(ctx context.Context, contactID, companyID, subject, body, errorMessage string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)
	log := EmailLogForSet{
		ID:                utils.GenerateUUID(),
		ContactID:         contactID,
		CompanyID:         companyID,
		Status:            StatusActive,
		EmailStage:        "ERROR",
		EmailSubject:      subject,
		EmailBody:         body,
		ErrorMessage:      errorMessage,
		AttachmentDetails: "{}",
		Metadata:          string(metadataJSON),
	}
	return InsertLog(ctx, log)
}
