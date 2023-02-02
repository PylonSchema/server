package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Email string
}

type AuthTokenClaims struct {
	TokenUUID string `json:"tid"`
	UserUUID  string `json:"uid"`
}

type RespInfo struct {
	Context *gin.Context
	Config  *oauth2.Config
	Token   *oauth2.Token
}

func ExchangeToken(c *gin.Context, code string, oauthConfig *oauth2.Config) (*oauth2.Token, error) {
	token, err := oauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Exhange error",
			"error":   err.Error(),
		})
	}
	return token, err
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GetLoginURL(state string, oauthConfig *oauth2.Config) string {
	return oauthConfig.AuthCodeURL(state)
}

func CheckState(c *gin.Context) error {
	cookie, err := c.Cookie("state")
	if err != nil {
		return err
	}
	session := sessions.Default(c)
	state := session.Get("state")
	if state == nil {
		return err
	}
	if state != cookie {
		return errors.New("state and cookie is not equal")
	}
	session.Delete("state")
	return nil
}

func (r *RespInfo) ReadBody(endpoint string) ([]byte, error) {
	var b []byte
	client := r.Config.Client(r.Context, r.Token)
	resp, err := client.Get(endpoint)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}

// func jwt() {

// }
