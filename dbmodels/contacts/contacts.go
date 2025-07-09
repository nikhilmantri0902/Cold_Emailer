package contacts

import (
	"cold_emailer/db"
	"context"
)

// ExistsByEmail checks if a contact exists by email, returns ID if found
func ExistsByEmail(ctx context.Context, email, status string) (string, error) {
	var id string
	query := `SELECT id FROM contacts WHERE email_id = $1 and status = $2 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, email, status)
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// InsertIfNotExists inserts a contact if it does not exist, returns the contact ID
func InsertIfNotExists(ctx context.Context, c ContactForSet) (string, error) {
	id, err := ExistsByEmail(ctx, c.EmailID, c.Status)
	if err == nil && id != "" {
		return id, nil // already exists
	}
	query := `INSERT INTO contacts (id, company_id, status, name, email_id, phone_number, linkedin_url, role, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err = db.GetDB().ExecContext(ctx, query, c.ID, c.CompanyID, c.Status, c.Name, c.EmailID, c.PhoneNumber, c.LinkedInURL, c.Role, c.Metadata)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}
