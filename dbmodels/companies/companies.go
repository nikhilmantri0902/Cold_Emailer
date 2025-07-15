package companies

import (
	"cold_emailer/db"
	"cold_emailer/utils"
	"context"
	"database/sql"
	"errors"
	"log"
)

// UpdateInput represents the fields to update for a company.
// Use pointer types so that nil means "do not update this field".
type UpdateInput struct {
	ID             string // required
	Status         *string
	ApolloID       *string
	Name           *string
	Website        *string
	Industry       *string
	SubIndustry    *string
	TechDetails    *string
	CompanyDetails *string
	Metadata       *string // JSON string
}

// Update updates the company with the given ID, only updating fields that are non-nil in input.
// Returns sql.ErrNoRows if the company does not exist.
func Update(ctx context.Context, input UpdateInput) error {
	if input.ID == "" {
		err := errors.New("id is not provided for update")
		log.Println("err:", err)
		return err
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.Status != nil {
		setClauses = append(setClauses, "status = $"+utils.Itoa(argIdx))
		args = append(args, *input.Status)
		argIdx++
	}
	if input.ApolloID != nil {
		setClauses = append(setClauses, "apollo_id = $"+utils.Itoa(argIdx))
		args = append(args, *input.ApolloID)
		argIdx++
	}
	if input.Name != nil {
		setClauses = append(setClauses, "name = $"+utils.Itoa(argIdx))
		args = append(args, *input.Name)
		argIdx++
	}
	if input.Website != nil {
		setClauses = append(setClauses, "website = $"+utils.Itoa(argIdx))
		args = append(args, *input.Website)
		argIdx++
	}
	if input.Industry != nil {
		setClauses = append(setClauses, "industry = $"+utils.Itoa(argIdx))
		args = append(args, *input.Industry)
		argIdx++
	}
	if input.SubIndustry != nil {
		setClauses = append(setClauses, "sub_industry = $"+utils.Itoa(argIdx))
		args = append(args, *input.SubIndustry)
		argIdx++
	}
	if input.TechDetails != nil {
		setClauses = append(setClauses, "tech_details = $"+utils.Itoa(argIdx))
		args = append(args, *input.TechDetails)
		argIdx++
	}
	if input.CompanyDetails != nil {
		setClauses = append(setClauses, "company_details = $"+utils.Itoa(argIdx))
		args = append(args, *input.CompanyDetails)
		argIdx++
	}
	if input.Metadata != nil {
		setClauses = append(setClauses, "metadata = $"+utils.Itoa(argIdx))
		args = append(args, *input.Metadata)
		argIdx++
	}

	if len(setClauses) == 0 {
		// Nothing to update
		return nil
	}

	query := "UPDATE companies SET " +
		utils.JoinClauses(setClauses, ", ") +
		" WHERE id = $" + utils.Itoa(argIdx)
	args = append(args, input.ID)

	res, err := db.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("err:", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("err:", err)
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetAll(ctx context.Context) (arr []Company, err error) {
	log.Println("inside get all companies")
	query := `SELECT id, status, coalesce(apollo_id, ''), 
						coalesce(name, '') as name, 
						coalesce(website, '') as website, 
						coalesce(industry, '') as industry, 
						coalesce(sub_industry, '') as sub_industry, 
						coalesce(tech_details, '') as tech_details, 
						coalesce(company_details, '') as company_details, 
						coalesce(metadata, '{}') as metadata
				FROM companies ORDER BY created_at DESC`

	rows, err := db.GetDB().QueryContext(ctx, query)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Company
		err := rows.Scan(
			&c.ID,
			&c.Status,
			&c.ApolloID,
			&c.Name,
			&c.Website,
			&c.Industry,
			&c.SubIndustry,
			&c.TechDetails,
			&c.CompanyDetails,
			&c.Metadata,
		)
		if err != nil {
			log.Println("err:", err)
			return nil, err
		}
		arr = append(arr, c)
	}

	if err := rows.Err(); err != nil {
		log.Println("err:", err)
		return nil, err
	}

	return arr, nil
}

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
