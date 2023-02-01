package github

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/devhoodit/sse-chat/auth"
	"github.com/devhoodit/sse-chat/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oAuthConfig *oauth2.Config

const (
	emailInfoEndpoint = "https://api.github.com/user/emails"
)

type Database interface {
	IsEmailUsed(email string) bool
	CreateUser(user *model.User) error
}

type Github struct {
	DB Database
}

type githubEmailInfo struct {
	Email    string
	Primary  bool
	Verified bool
}

func init() {
	oAuthConfig = &oauth2.Config{
		ClientID:     "03310852bd9891db5f0e",
		ClientSecret: "e2989c0dbb1896a097882778fb05ba5f9fc02e4a",
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func (g *Github) Login(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/auth",
		MaxAge: 900,
	})
	state := auth.RandToken()
	session.Set("state", state)
	session.Save()
	c.SetCookie("state", state, 900, "/auth", "localhost", true, false)
	c.Redirect(http.StatusFound, auth.GetLoginURL(state, oAuthConfig))
}

func (g *Github) Callback(c *gin.Context) {

	cookie, err := c.Cookie("state")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No state cookie",
		})
		return
	}

	session := sessions.Default(c)
	state := session.Get("state")
	if state == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "session is nil",
		})
		return
	}

	if state != cookie {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": "AccessDenied",
		})
		return
	}

	session.Delete("state")

	token, err := oAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Exchange error",
			"error":   err.Error(),
		})
		return
	}

	client := oAuthConfig.Client(c, token)
	userInfoResp, err := client.Get(emailInfoEndpoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "code resp error",
		})
		return
	}

	defer userInfoResp.Body.Close()
	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "read resp body error",
		})
		return
	}

	var infos []githubEmailInfo

	err = json.Unmarshal(userInfo, &infos)
	if err != nil {
		panic(err)
	}

	var email string = ""

	for _, info := range infos {
		if info.Verified {
			email = info.Email
			break
		}
	}
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No vaild email",
		})
	}

	// extraction email
	if g.DB.IsEmailUsed(email) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "email is already used",
		})
		return
	}

	err = g.createUser(email, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error, try again",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok, user successfully created",
	})
}

func (g *Github) createUser(email string, token *oauth2.Token) error {
	privateUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	publicUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user := model.User{
		Username:    email,
		AccountType: 1, // static, account type is social
		UUID:        publicUUID,
		SecretUUID:  privateUUID,
		Email:       email,
	}

	err = g.DB.CreateUser(&user)
	if err != nil {
		return err
	}

	return nil
}
