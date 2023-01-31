package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Email string
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
