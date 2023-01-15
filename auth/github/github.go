package github

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var OAuthConfig *oauth2.Config

const (
	UserInfoEndpoint = "https://api.github.com/user/emails"
)

type GithubEmailInfo struct {
	Email    string
	Primary  bool
	Verified bool
}

func init() {
	OAuthConfig = &oauth2.Config{
		ClientID:     "03310852bd9891db5f0e",
		ClientSecret: "e2989c0dbb1896a097882778fb05ba5f9fc02e4a",
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email"},
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
	session.Save()
	c.SetCookie("state", state, 900, "/auth", "localhost", true, false)
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

func Authenticate(c *gin.Context) {

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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Wrong state",
			"state":   state,
			"cookie":  cookie,
		})
		return
	}

	session.Delete("state")

	token, err := OAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Exchange error",
			"error":   err.Error(),
		})
		return
	}

	client := OAuthConfig.Client(c, token)
	userInfoResp, err := client.Get(UserInfoEndpoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "code resp error",
		})
		return
	}

	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "read resp body error",
		})
		return
	}

	var infos []GithubEmailInfo

	err = json.Unmarshal(userInfo, &infos)
	if err != nil {
		panic(err)
	}

	var email string = ""

	for _, info := range infos {
		if !info.Primary {
			continue
		}
		if !info.Verified {
			continue
		}
		email = info.Email
	}
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No vaild email",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
