package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	JwtAuth *JwtAuth
	d       AuthDatabase
}

type AuthDatabase interface {
}

func New(jwtAuth *JwtAuth, authDatabase AuthDatabase) *AuthAPI {
	return &AuthAPI{
		JwtAuth: jwtAuth,
		d:       authDatabase,
	}
}

func (a *AuthAPI) BlacklistHandler(c *gin.Context) {
	token, err := c.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "no token cookie",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no valid token",
		})
		return
	}

	// parse cookie
	tokenString := token.Value
	claims := &AuthTokenClaims{}
	jwtToken, err := a.JwtAuth.ParseToken(c, claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	// check token is valid
	if !jwtToken.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "token is invalid",
		})
		return
	}
	err = a.JwtAuth.Store.SetBlacklist(tokenString, time.Until(claims.ExpiresAt.Time))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": time.Until(claims.ExpiresAt.Time),
	})
}

func (a *AuthAPI) GetTokenHandler(c *gin.Context) {

}
