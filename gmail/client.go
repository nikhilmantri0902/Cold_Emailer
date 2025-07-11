package gmail

import (
	"cold_emailer/constants"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var gmailScopes = []string{
	"https://www.googleapis.com/auth/gmail.send",
	"https://www.googleapis.com/auth/gmail.readonly",
}

// GetOAuth2Config returns the OAuth2 config for Gmail
func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     constants.GMAIL_CLIENT_ID,
		ClientSecret: constants.GMAIL_CLIENT_SECRET,
		RedirectURL:  constants.GMAIL_REDIRECT_URI,
		Scopes:       gmailScopes,
		Endpoint:     google.Endpoint,
	}
}

// GetAuthURL returns the URL to redirect the user to for Gmail OAuth2 consent
func GetAuthURL(state string) string {
	config := GetOAuth2Config()
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCode exchanges the code for a token (access + refresh)
func ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	config := GetOAuth2Config()
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Println("error:", err)
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return tok, nil
}

// SendSingleEmail sends an email using the Gmail API and the given OAuth2 token
func SendSingleEmail(ctx context.Context, accessToken, from, to, subject, body string) error {
	// Create the message
	msg := &mail.Message{
		Header: map[string][]string{
			"From":    {from},
			"To":      {to},
			"Subject": {subject},
		},
		Body: strings.NewReader(body),
	}

	var msgBuilder strings.Builder
	for k, v := range msg.Header {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(body)

	raw := base64.URLEncoding.EncodeToString([]byte(msgBuilder.String()))

	// Create Gmail service
	svc, err := gmail.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Hour), // not used for static token
	})))
	if err != nil {
		log.Println("error:", err)
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	// Send the message
	_, err = svc.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Do()
	if err != nil {
		log.Println("error:", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendEmailWithAttachment sends an email with resume attachment using Gmail API
func SendEmailWithAttachment(ctx context.Context, accessToken, from, to, subject, body, resumePath string) error {
	// Read resume file
	resumeData, err := os.ReadFile(resumePath)
	if err != nil {
		return fmt.Errorf("failed to read resume file: %w", err)
	}

	// Get file extension for MIME type
	ext := strings.ToLower(filepath.Ext(resumePath))
	var mimeType string
	switch ext {
	case ".pdf":
		mimeType = "application/pdf"
	case ".doc":
		mimeType = "application/msword"
	case ".docx":
		mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		mimeType = "application/octet-stream"
	}

	// Create email with attachment
	boundary := "boundary123"
	var emailBuilder strings.Builder

	// Email headers
	emailBuilder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	emailBuilder.WriteString(fmt.Sprintf("To: %s\r\n", to))
	emailBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	emailBuilder.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\n"))
	emailBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	emailBuilder.WriteString("\r\n")

	// Email body
	emailBuilder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	emailBuilder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	emailBuilder.WriteString("\r\n")
	emailBuilder.WriteString(body)
	emailBuilder.WriteString("\r\n")

	// Resume attachment
	emailBuilder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	emailBuilder.WriteString(fmt.Sprintf("Content-Type: %s; name=\"resume%s\"\r\n", mimeType, ext))
	emailBuilder.WriteString("Content-Transfer-Encoding: base64\r\n")
	emailBuilder.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"resume%s\"\r\n", ext))
	emailBuilder.WriteString("\r\n")
	emailBuilder.WriteString(base64.StdEncoding.EncodeToString(resumeData))
	emailBuilder.WriteString("\r\n")

	// End boundary
	emailBuilder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	raw := base64.URLEncoding.EncodeToString([]byte(emailBuilder.String()))

	// Create Gmail service
	svc, err := gmail.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Hour),
	})))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	// Send the message
	_, err = svc.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
