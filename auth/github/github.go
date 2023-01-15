package github

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var OAuthConfig *oauth2.Config

func init() {
	OAuthConfig = &oauth2.Config{
		ClientID:     "03310852bd9891db5f0e",
		ClientSecret: "e2989c0dbb1896a097882778fb05ba5f9fc02e4a",
		RedirectURL:  "https://localhost:8080/auth/github/callback",
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
}

func RenderAuthView(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/auth",
		MaxAge: 900,
	})
	state := RandToken()
	session.Set("state", state)
	c.SetCookie("state", state, 900, "/auth", "localhost", true, true)
	c.Redirect(http.StatusFound, getLoginURL(state))
}

func getLoginURL(state string) string {
	return OAuthConfig.AuthCodeURL(state)
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
