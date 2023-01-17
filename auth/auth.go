package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

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
