package email_logs

type EmailLog struct {
	ID                string `db:"id"`
	ContactID         string `db:"contact_id"`
	CompanyID         string `db:"company_id"`
	Status            string `db:"status"`
	EmailStage        string `db:"email_stage"`
	EmailSubject      string `db:"email_subject"`
	EmailBody         string `db:"email_body"`
	AttachmentDetails string `db:"attachment_details"`
	ErrorMessage      string `db:"error_message"`
	Metadata          string `db:"metadata"`
}

type EmailLogForSet struct {
	ID                string
	ContactID         string
	CompanyID         string
	Status            string
	EmailStage        string
	EmailSubject      string
	EmailBody         string
	AttachmentDetails string
	ErrorMessage      string
	Metadata          string
}
