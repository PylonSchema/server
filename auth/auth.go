package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	JwtAuth *JwtAuth
}

func (a *Auth) Blacklist(c *gin.Context) {
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
