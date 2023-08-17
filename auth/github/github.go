package github

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	emailInfoEndpoint = "https://api.github.com/user/emails"
	userInfoEndpoint  = "https://api.github.com/user"
)

type Database interface {
	IsEmailUsed(email string) (bool, error)
	GetUserFromSocialByEmail(email string, socialType int) (*model.User, error)
	CreateUser(user *model.User) error
	CreateSocial(social *model.Social) error
}

type Github struct {
	DB          Database
	JwtAuth     *auth.JwtAuth
	OAuthConfig *oauth2.Config
}

type githubUserInfo struct {
	Login string
}

type githubEmailInfo struct {
	Email    string
	Primary  bool
	Verified bool
}

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

func (g *Github) createUser(username string, userId string, email string, token *oauth2.Token) error {

	publicUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user := model.User{
		Username:    username,
		AccountType: 1, // static, account type is social
		UUID:        publicUUID,
		Email:       email,
	}
	err = g.DB.CreateUser(&user)
	if err != nil {
		return err
	}
	social := model.Social{
		SocialType:   1, // static account type is github,
		SocialId:     userId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken, // this will be nil, github has no refresh token
	}
	err = g.DB.CreateSocial(&social)
	if err != nil {
		return err
	}

	return nil
}
