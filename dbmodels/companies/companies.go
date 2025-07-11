package companies

import (
	"cold_emailer/db"
	"context"
	"database/sql"
	"log"
)

// ExistsByApolloIDBool checks if a company exists by apollo_id and status
// Returns (true, nil) if exists, (false, nil) if not, (false, err) for other errors
func ExistsByApolloID(ctx context.Context, apolloID, status string) (bool, error) {
	var id string
	query := `SELECT id FROM companies WHERE apollo_id = $1 AND status = $2 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, apolloID, status)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		log.Println("error:", err)
		return false, err
	}
	return true, nil
}

// ExistsByNameWebsite checks if a company exists by name, website, and status
// Returns (true, nil) if exists, (false, nil) if not, (false, err) for other errors
func ExistsByNameWebsite(ctx context.Context, name, website, status string) (bool, error) {
	var id string
	query := `SELECT id FROM companies WHERE name = $1 AND website = $2 AND status = $3 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, name, website, status)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		log.Println("error:", err)
		return false, err
	}
	return true, nil
}

// InsertIfNotExists inserts a company if it does not exist, returns the company ID
func InsertIfNotExists(ctx context.Context, c CompanyForSet) (string, error) {
	exists, err := ExistsByNameWebsite(ctx, c.Name, c.Website, c.Status)
	if err == nil && exists {
		return "", nil // already exists
	}
	query := `INSERT INTO companies (id, apollo_id, status, name, website, industry, sub_industry, tech_details, company_details, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = db.GetDB().ExecContext(ctx, query, c.ID, c.ApolloID, c.Status, c.Name, c.Website, c.Industry, c.SubIndustry, c.TechDetails, c.CompanyDetails, c.Metadata)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func DeleteByID(ctx context.Context, id string) (err error) {
	query := `DELETE from companies where id = $1;`
	_, err = db.GetDB().ExecContext(ctx, query, id)
	if err != nil {
		log.Println("err:", err)
		return err
	}
	return nil
}
