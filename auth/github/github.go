package github

import (
	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
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
