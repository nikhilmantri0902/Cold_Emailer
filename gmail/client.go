package gmail

import (
	"cold_emailer/constants"
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
