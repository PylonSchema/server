package github

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (g *Github) LoginHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/auth",
		MaxAge: 900,
	})
	state := auth.RandToken()
	session.Set("state", state)
	session.Save()
	c.SetCookie("state", state, 900, "/auth", "localhost", true, false)
	c.Redirect(http.StatusFound, auth.GetLoginURL(state, g.OAuthConfig))
}
