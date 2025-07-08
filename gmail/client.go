package gmail

import (
	"cold_emailer/constants"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
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
