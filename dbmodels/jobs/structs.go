package jobs

import "time"

type Job struct {
	ID          string     `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Status      string     `db:"status"`
	Name        string     `db:"name"`
	Message     string     `db:"message"`
	Metadata    string     `db:"metadata"`
}

type JobForSet struct {
	ID          string
	CompletedAt *time.Time
	Status      string
	Name        string
	Message     string
	Metadata    string
}

// UpdateInput represents the fields to update for a job.
// Use pointer types so that nil means "do not update this field".
type UpdateInput struct {
	ID          string // required
	CompletedAt **time.Time
	Status      *string
	Name        *string
	Message     *string
	Metadata    *string // JSON string
}
