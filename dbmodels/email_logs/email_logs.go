package email_logs

import (
	"cold_emailer/db"
	"cold_emailer/utils"
	"context"
	"encoding/json"
)

// EmailLogWithDetails represents an email log with company and contact information
type EmailLogWithDetails struct {
	ID                string `json:"id"`
	ContactID         string `json:"contact_id"`
	CompanyID         string `json:"company_id"`
	Status            string `json:"status"`
	EmailStage        string `json:"email_stage"`
	EmailSubject      string `json:"email_subject"`
	EmailBody         string `json:"email_body"`
	AttachmentDetails string `json:"attachment_details"`
	ErrorMessage      string `json:"error_message"`
	Metadata          string `json:"metadata"`
	CreatedAt         string `json:"created_at"`
	// Company details
	CompanyName     string `json:"company_name"`
	CompanyWebsite  string `json:"company_website"`
	CompanyIndustry string `json:"company_industry"`
	// Contact details
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactRole  string `json:"contact_role"`
}

// InsertLog inserts a new email log entry
func InsertLog(ctx context.Context, l EmailLogForSet) error {
	query := `INSERT INTO email_logs (id, contact_id, company_id, status, email_stage, email_subject, email_body, attachment_details, error_message, metadata, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, now())`
	_, err := db.GetDB().ExecContext(ctx, query, l.ID, l.ContactID, l.CompanyID, l.Status, l.EmailStage, l.EmailSubject, l.EmailBody, l.AttachmentDetails, l.ErrorMessage, l.Metadata)
	return err
}

// GetEmailLogsWithPagination fetches email logs with pagination and optional stage filter
func GetEmailLogsWithPagination(ctx context.Context, offset, limit int, stage string) ([]EmailLogWithDetails, int, error) {
	// Build the query with optional stage filter
	baseQuery := `
		SELECT 
			el.id, el.contact_id, el.company_id, el.status, el.email_stage,
			el.email_subject, el.email_body, el.attachment_details, el.error_message,
			el.metadata, el.created_at,
			comp.name as company_name, comp.website as company_website, comp.industry as company_industry,
			c.name as contact_name, c.email_id as contact_email, c.role as contact_role
		FROM email_logs el
		LEFT JOIN companies comp ON el.company_id = comp.id
		LEFT JOIN contacts c ON el.contact_id = c.id
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM email_logs el
	`

	var whereClause string
	var args []interface{}
	argIndex := 1

	if stage != "" {
		whereClause = " WHERE el.email_stage = $" + utils.Itoa(argIndex)
		args = append(args, stage)
		argIndex++
	}

	// Execute count query
	var total int
	err := db.GetDB().QueryRowContext(ctx, countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("err:")
		return nil, 0, err
	}

	// Execute main query with pagination
	finalQuery := baseQuery + whereClause + " ORDER BY el.created_at DESC LIMIT $" + utils.Itoa(argIndex) + " OFFSET $" + utils.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.GetDB().QueryContext(ctx, finalQuery, args...)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("err:")
		return nil, 0, err
	}
	defer rows.Close()

	var logs []EmailLogWithDetails
	for rows.Next() {
		var log EmailLogWithDetails
		err := rows.Scan(
			&log.ID, &log.ContactID, &log.CompanyID, &log.Status, &log.EmailStage,
			&log.EmailSubject, &log.EmailBody, &log.AttachmentDetails, &log.ErrorMessage,
			&log.Metadata, &log.CreatedAt,
			&log.CompanyName, &log.CompanyWebsite, &log.CompanyIndustry,
			&log.ContactName, &log.ContactEmail, &log.ContactRole,
		)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
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
