package email_logs

import (
	"cold_emailer/db"
	"context"
)

// InsertLog inserts a new email log entry
func InsertLog(ctx context.Context, l EmailLogForSet) error {
	query := `INSERT INTO email_logs (id, contact_id, company_id, status, email_stage, email_subject, email_body, attachment_details, error_message, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := db.GetDB().ExecContext(ctx, query, l.ID, l.ContactID, l.CompanyID, l.Status, l.EmailStage, l.EmailSubject, l.EmailBody, l.AttachmentDetails, l.ErrorMessage, l.Metadata)
	return err
}
