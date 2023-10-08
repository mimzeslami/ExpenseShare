package util

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOAuthConfig() (*oauth2.Config, error) {
	config, err := LoadConfig(".")
	if err != nil {
		return nil, err
	}
	outhConfig := &oauth2.Config{
		ClientID:     config.GoogleOAuthClientID,
		ClientSecret: config.GoogleOAuthSecret,
		RedirectURL:  config.GoogleOAuthRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return outhConfig, nil
}
