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
