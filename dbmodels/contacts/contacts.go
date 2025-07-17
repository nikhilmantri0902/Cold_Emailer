package contacts

import (
	"cold_emailer/db"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// ExistsByApolloID checks if a contact exists by apollo_id and status
// Returns (true, nil) if exists, (false, nil) if not, (false, err) for other errors
func ExistsByApolloID(ctx context.Context, apolloID, status string) (bool, error) {
	var id string
	query := `SELECT id FROM contacts WHERE apollo_id = $1 AND status = $2 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, apolloID, status)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ExistsByEmail checks if a contact exists by email and status
// Returns (true, nil) if exists, (false, nil) if not, (false, err) for other errors
func ExistsByEmail(ctx context.Context, email, status string) (bool, error) {
	var id string
	query := `SELECT id FROM contacts WHERE email_id = $1 AND status = $2 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, email, status)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// InsertIfNotExists inserts a contact if it does not exist, returns the contact ID
func InsertIfNotExists(ctx context.Context, c ContactForSet) (string, error) {
	id, err := ExistsByEmail(ctx, c.EmailID, c.Status)
	if err == nil && id {
		return "", nil // already exists
	}
	query := `INSERT INTO contacts (id, apollo_id, company_id, status, name, email_id, phone_number, linkedin_url, role, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = db.GetDB().ExecContext(ctx, query, c.ID, c.ApolloID, c.CompanyID, c.Status, c.Name, c.EmailID, c.PhoneNumber, c.LinkedInURL, c.Role, c.Metadata)
	if err != nil {
		log.Println("err:", err)
		return "", err
	}
	return c.ID, nil
}

// GetContactsWithCompanyInfo fetches contacts with their company information, excluding those with SENT emails
func GetContactsWithCompanyInfo(ctx context.Context, count int, status, orderBy string) ([]ContactWithCompany, error) {
	if orderBy == "" {
		orderBy = "c.created_at DESC"
	}

	query := `
		SELECT 
			c.id as contact_id,
			c.name as contact_name,
			c.email_id as contact_email,
			c.role as contact_role,
			c.linkedin_url as contact_linkedin,
			c.phone_number as contact_phone,
			c.status as contact_status,
			comp.id as company_id,
			comp.name as company_name,
			comp.website as company_website,
			comp.industry as company_industry,
			comp.tech_details as tech_details,
			comp.company_details as company_details
		FROM contacts c
		JOIN companies comp ON c.company_id = comp.id
		LEFT JOIN email_logs el ON c.id = el.contact_id AND el.email_stage = 'SENT'
		WHERE c.status = $1
		AND el.id IS NULL
		ORDER BY ` + orderBy + `
		LIMIT $2
	`

	rows, err := db.GetDB().QueryContext(ctx, query, status, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactWithCompany
	for rows.Next() {
		var contact ContactWithCompany
		err := rows.Scan(
			&contact.ContactID,
			&contact.ContactName,
			&contact.ContactEmail,
			&contact.ContactRole,
			&contact.ContactLinkedIn,
			&contact.ContactPhone,
			&contact.ContactStatus,
			&contact.CompanyID,
			&contact.CompanyName,
			&contact.CompanyWebsite,
			&contact.CompanyIndustry,
			&contact.CompanyTech,
			&contact.CompanyDetails,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetFollowUpCandidates fetches contacts whose latest email log is SENT, ACTIVE, and sent_at > 3 days ago
func GetFollowUpCandidates(ctx context.Context, daysPastFirstEmail int, limit int, status string) ([]FollowUpCandidate, error) {
	query :=
		`SELECT
			el.contact_id,
			el.company_id,
			c.name as contact_name,
			c.email_id as contact_email,
			c.role as contact_role,
			c.linkedin_url as contact_linkedin,
			c.phone_number as contact_phone,
			c.status as contact_status,
			comp.name as company_name,
			comp.website as company_website,
			comp.industry as company_industry,
			comp.tech_details as tech_details,
			comp.company_details as company_details,
			el.email_subject,
			el.email_body,
			el.metadata
		FROM email_logs el
		JOIN (
			SELECT contact_id, MAX(created_at) as max_created_at
			FROM email_logs
			GROUP BY contact_id
		) latest ON el.contact_id = latest.contact_id AND el.created_at = latest.max_created_at
		JOIN contacts c ON el.contact_id = c.id
		JOIN companies comp ON el.company_id = comp.id
		WHERE el.status = $1
		  AND el.email_stage = $2
		  AND (el.metadata->>'sent_at')::timestamp < (NOW() - INTERVAL '%d days')
		ORDER BY (el.metadata->>'sent_at')::timestamp ASC
		LIMIT $3;
	`

	query = fmt.Sprintf(query, daysPastFirstEmail)
	rows, err := db.GetDB().QueryContext(ctx, query, status, "SENT", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []FollowUpCandidate
	for rows.Next() {
		var c FollowUpCandidate
		err := rows.Scan(
			&c.ContactID,
			&c.CompanyID,
			&c.ContactName,
			&c.ContactEmail,
			&c.ContactRole,
			&c.ContactLinkedIn,
			&c.ContactPhone,
			&c.ContactStatus,
			&c.CompanyName,
			&c.CompanyWebsite,
			&c.CompanyIndustry,
			&c.CompanyTech,
			&c.CompanyDetails,
			&c.EmailSubject,
			&c.EmailBody,
			&c.Metadata,
		)
		if err != nil {
			log.Println("err:", err)
			return nil, err
		}
		candidates = append(candidates, c)
	}
	return candidates, nil
}
