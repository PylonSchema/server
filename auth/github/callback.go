package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (g *Github) CallbackHandler(c *gin.Context) {

	err := auth.CheckState(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, auth.AccessDenied)
		return
	}

	token, err := g.OAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, auth.AccessDenied)
		return
	}

	userId, err := g.getUserId(c, token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, auth.AccessDenied)
		return
	}

	email, err := g.getUserEmail(c, token)
	if err == auth.ErrNoValidEmail {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "No vaild email, have no certified email",
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, auth.AccessDenied)
		return
	}

	isEmailUsed, err := g.DB.IsEmailUsed(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
		fmt.Println("github auth email is used")
		return
	}
	if !isEmailUsed {
		// Create User Flow
		err = g.createUser(userId, userId, email, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
			fmt.Println("github auth create user error")
			return
		}
	}

	// validate email is on this platform
	user, err := g.DB.GetUserFromSocialByEmail(email, 1) // github social type is 1
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
		fmt.Println("github auth get user by email error")
		return
	}

	jp := auth.JwtPayload{
		UserUUID: user.UUID,
		Username: user.Username,
	}
	jwtTokenString, err := g.JwtAuth.GenerateJWT(&jp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
		fmt.Println("github auth generate jwt token error")
		return
	}
	c.SetCookie("token", jwtTokenString, 60*60*24*90, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
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
		return "", auth.ErrNoValidEmail
	}
	return email, nil
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
