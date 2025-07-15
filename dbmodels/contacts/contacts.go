package contacts

import (
	"cold_emailer/db"
	"context"
	"database/sql"
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

// ContactWithCompany represents a contact with company information
type ContactWithCompany struct {
	ContactID       string `db:"contact_id"`
	ContactName     string `db:"contact_name"`
	ContactEmail    string `db:"contact_email"`
	ContactRole     string `db:"contact_role"`
	ContactLinkedIn string `db:"contact_linkedin"`
	ContactPhone    string `db:"contact_phone"`
	ContactStatus   string `db:"contact_status"`
	CompanyID       string `db:"company_id"`
	CompanyName     string `db:"company_name"`
	CompanyWebsite  string `db:"company_website"`
	CompanyIndustry string `db:"company_industry"`
	CompanyTech     string `db:"tech_details"`
	CompanyDetails  string `db:"company_details"`
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
