package jobs

import (
	"cold_emailer/db"
	"cold_emailer/utils"
	"context"
	"database/sql"
	"errors"
	"log"
)

// Insert inserts a new job into the database
func Insert(ctx context.Context, job JobForSet) (string, error) {
	if job.ID == "" {
		return "", errors.New("id is required for insert")
	}

	metadata := "{}"
	if job.Metadata != "" {
		metadata = job.Metadata
	}

	query := `INSERT INTO jobs (id, status, name, message, completed_at, metadata) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.GetDB().ExecContext(ctx, query, job.ID, job.Status, job.Name, job.Message, job.CompletedAt, metadata)
	if err != nil {
		log.Println("err:", err)
		return "", err
	}
	return job.ID, nil
}

// Update updates the job with the given ID, only updating fields that are non-nil in input.
// Returns sql.ErrNoRows if the job does not exist.
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
	if input.Name != nil {
		setClauses = append(setClauses, "name = $"+utils.Itoa(argIdx))
		args = append(args, *input.Name)
		argIdx++
	}
	if input.Message != nil {
		setClauses = append(setClauses, "message = $"+utils.Itoa(argIdx))
		args = append(args, *input.Message)
		argIdx++
	}
	if input.CompletedAt != nil {
		setClauses = append(setClauses, "completed_at = $"+utils.Itoa(argIdx))
		args = append(args, *input.CompletedAt)
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

	query := "UPDATE jobs SET " +
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

// GetByID retrieves a job by its ID, optionally filtered by status
func GetByID(ctx context.Context, id string, status string) (Job, error) {
	var job Job
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT 
			id, 
			created_at, 
			completed_at,
			COALESCE(status, '') as status,
			COALESCE(name, '') as name,
			COALESCE(message, '') as message,
			COALESCE(metadata, '{}') as metadata
		FROM jobs 
		WHERE id = $1 AND status = $2`
		args = []interface{}{id, status}
	} else {
		query = `SELECT 
			id, 
			created_at, 
			completed_at,
			COALESCE(status, '') as status,
			COALESCE(name, '') as name,
			COALESCE(message, '') as message,
			COALESCE(metadata, '{}') as metadata
		FROM jobs 
		WHERE id = $1`
		args = []interface{}{id}
	}

	row := db.GetDB().QueryRowContext(ctx, query, args...)
	err := row.Scan(
		&job.ID,
		&job.CreatedAt,
		&job.CompletedAt,
		&job.Status,
		&job.Name,
		&job.Message,
		&job.Metadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			if status != "" {
				log.Printf("Job with ID %s and status %s not found", id, status)
			} else {
				log.Printf("Job with ID %s not found", id)
			}
		} else {
			log.Println("err:", err)
		}
		return job, err
	}
	return job, nil
}

// GetByName retrieves a job by its name, optionally filtered by status
func GetByName(ctx context.Context, name string, status string) (Job, error) {
	var job Job
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT 
			id, 
			created_at, 
			completed_at,
			COALESCE(status, '') as status,
			COALESCE(name, '') as name,
			COALESCE(message, '') as message,
			COALESCE(metadata, '{}') as metadata
		FROM jobs 
		WHERE name = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1`
		args = []interface{}{name, status}
	} else {
		query = `SELECT 
			id, 
			created_at, 
			completed_at,
			COALESCE(status, '') as status,
			COALESCE(name, '') as name,
			COALESCE(message, '') as message,
			COALESCE(metadata, '{}') as metadata
		FROM jobs 
		WHERE name = $1
		ORDER BY created_at DESC
		LIMIT 1`
		args = []interface{}{name}
	}

	row := db.GetDB().QueryRowContext(ctx, query, args...)
	err := row.Scan(
		&job.ID,
		&job.CreatedAt,
		&job.CompletedAt,
		&job.Status,
		&job.Name,
		&job.Message,
		&job.Metadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			if status != "" {
				log.Printf("Job with name %s and status %s not found", name, status)
			} else {
				log.Printf("Job with name %s not found", name)
			}
		} else {
			log.Println("err:", err)
		}
		return job, err
	}
	return job, nil
}
