package github

import (
	"encoding/json"
	"net/http"

	"github.com/devhoodit/sse-chat/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (g *Github) Callback(c *gin.Context) {
	err := auth.CheckState(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, AccessDenied)
		return
	}

	token, err := g.OAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, AccessDenied)
		return
	}

	userId, err := g.getUserId(c, token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, AccessDenied)
		return
	}

	email, err := g.getUserEmail(c, token)
	if err == auth.ErrNoVaildEmail {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "No vaild email, have no certified email",
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, AccessDenied)
		return
	}

	if !g.DB.IsEmailUsed(email) {
		// Create User Flow
		err = g.createUser(userId, userId, email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "server error, try again",
			})
			return
		}
	}

	// login flow
	c.JSON(http.StatusOK, gin.H{
		"message": "ok, user successfully created",
	})
}

func (g *Github) getUserId(c *gin.Context, token *oauth2.Token) (string, error) {
	client := auth.RespInfo{
		Context: c,
		Config:  g.OAuthConfig,
		Token:   token,
	}

	userInfo, err := client.ReadBody(userInfoEndpoint)
	if err != nil {
		return "", err
	}

	var info githubUserInfo
	err = json.Unmarshal(userInfo, &info)
	if err != nil {
		return "", err
	}

	return info.Login, nil
}

func (g *Github) getUserEmail(c *gin.Context, token *oauth2.Token) (string, error) {
	client := auth.RespInfo{
		Context: c,
		Config:  g.OAuthConfig,
		Token:   token,
	}

	emailInfo, err := client.ReadBody(emailInfoEndpoint)
	if err != nil {
		return "", err
	}

	var infos []githubEmailInfo
	err = json.Unmarshal(emailInfo, &infos)
	if err != nil {
		return "", err
	}

	var email string = ""
	for _, info := range infos {
		if info.Verified {
			email = info.Email
			break
		}
	}
	if email == "" {
		return "", auth.ErrNoVaildEmail
	}
	return email, nil
}
