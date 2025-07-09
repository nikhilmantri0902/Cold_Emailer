package companies

import (
	"cold_emailer/db"
	"context"
)

// ExistsByNameWebsite checks if a company exists by name and website, returns ID if found
func ExistsByNameWebsite(ctx context.Context, name, website, status string) (string, error) {
	var id string
	query := `SELECT id FROM companies WHERE name = $1 AND website = $2 AND status = $3 LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, name, website, status)
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// InsertIfNotExists inserts a company if it does not exist, returns the company ID
func InsertIfNotExists(ctx context.Context, c CompanyForSet) (string, error) {
	id, err := ExistsByNameWebsite(ctx, c.Name, c.Website, c.Status)
	if err == nil && id != "" {
		return id, nil // already exists
	}
	query := `INSERT INTO companies (id, status, name, website, industry, sub_industry, tech_details, company_details, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err = db.GetDB().ExecContext(ctx, query, c.ID, c.Status, c.Name, c.Website, c.Industry, c.SubIndustry, c.TechDetails, c.CompanyDetails, c.Metadata)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}
