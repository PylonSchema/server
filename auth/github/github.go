package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var OAuthConfig *oauth2.Config

func init() {
	const (
		clientID     = "03310852bd9891db5f0e"
		clientSecret = "e2989c0dbb1896a097882778fb05ba5f9fc02e4a"
		redirectURL  = "https://localhost:8080"
		ScopeEmail   = ""
	)

	OAuthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://localhost:8080",
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
}

func getLoginURL(state string) string {
	return OAuthConfig.AuthCodeURL(state)
}
