package gmailtokens

import (
	"cold_emailer/db"
	"context"
	"log"
	"time"
)

type GmailTokenForSet struct {
	ID           string
	EmailID      string
	AccessToken  string
	RefreshToken string
	Expiry       string // RFC3339 or timestamp string
}

type GmailToken struct {
	ID           string
	EmailID      string
	AccessToken  string
	RefreshToken string
	Expiry       string
	CreatedAt    string
}

// CheckExpiry returns true if expired.
func (g *GmailToken) CheckExpiry() bool {
	expiryTime, err := time.Parse(time.RFC3339, g.Expiry)
	if err != nil {
		log.Println("err:", err)
		return true // If parsing fails, treat as expired
	}
	return time.Now().After(expiryTime)
}

func CreateGmailToken(ctx context.Context, t GmailTokenForSet) error {
	query := `INSERT INTO gmail_tokens (
		id, email_id, access_token, refresh_token, expiry
	) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.GetDB().ExecContext(ctx, query,
		t.ID,
		t.EmailID,
		t.AccessToken,
		t.RefreshToken,
		t.Expiry,
	)
	if err != nil {
		log.Println("error:", err)
	}
	return err
}

func GetLatestToken(ctx context.Context) (GmailToken, error) {
	var t GmailToken
	query := `SELECT id, email_id, access_token, refresh_token, expiry, created_at FROM gmail_tokens ORDER BY created_at DESC LIMIT 1`
	row := db.GetDB().QueryRowContext(ctx, query)
	err := row.Scan(&t.ID, &t.EmailID, &t.AccessToken, &t.RefreshToken, &t.Expiry, &t.CreatedAt)
	if err != nil {
		log.Println("error:", err)
	}
	return t, err
}
