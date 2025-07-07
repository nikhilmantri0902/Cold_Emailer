package profileinfo

import (
	"cold_emailer/db"
	"context"
	"log"
)

func Create(ctx context.Context, info StructForSet) error {

	metadata := "{}"
	if info.Metadata != "" {
		metadata = info.Metadata
	}
	query := `INSERT INTO profile_info (
		id, status, name, email, phone, linkedin_url, experience, skills, summary, resume_path, metadata
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.GetDB().ExecContext(ctx, query,
		info.ID,
		info.Status,
		info.Name,
		info.Email,
		info.Phone,
		info.LinkedInURL,
		info.Experience,
		info.Skills,
		info.Summary,
		info.ResumePath,
		metadata,
	)
	if err != nil {
		log.Println("error:", err)
	}
	return err
}

func GetLatestActive(ctx context.Context) (Struct, error) {
	var result Struct
	query := `SELECT 
			id, 
			created_at, 
			COALESCE(status, '') as status, 
			name, 
			email, 
			COALESCE(phone, '') as phone, 
			COALESCE(linkedin_url, '') as linkedin_url, 
			COALESCE(experience, '') as experience, 
			COALESCE(skills, '') as skills, 
			COALESCE(summary, '') as summary, 
			COALESCE(resume_path, '') as resume_path,
			COALESCE(metadata, '{}') as metadata
		FROM profile_info
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query, StatusActive)
	err := row.Scan(
		&result.ID,
		&result.CreatedAt,
		&result.Status,
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.LinkedInURL,
		&result.Experience,
		&result.Skills,
		&result.Summary,
		&result.ResumePath,
		&result.Metadata,
	)
	return result, err
}
