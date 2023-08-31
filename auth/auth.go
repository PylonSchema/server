package auth

import (
	"net/http"
	"time"

	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthAPI struct {
	JwtAuth *JwtAuth
	d       AuthDatabase
}

type UserToken struct {
	ExpireAt   time.Time `json:"expire_at"`
	Type       int       `json:"type"`
	DeviceName string    `json:"device_name"`
}

type AuthDatabase interface {
	GetAllUserToken(uuid uuid.UUID) (*[]model.UserTokenPair, error)
}

func New(jwtAuth *JwtAuth, authDatabase AuthDatabase) *AuthAPI {
	return &AuthAPI{
		JwtAuth: jwtAuth,
		d:       authDatabase,
	}
}

func (a *AuthAPI) BlacklistHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*AuthTokenClaims)
	tokenString := c.Request.Header.Get("X-Pylon-Token")
	err := a.JwtAuth.Store.SetBlacklist(tokenString, time.Until(claims.ExpiresAt.Time))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": time.Until(claims.ExpiresAt.Time),
	})
}

func (a *AuthAPI) GetTokenHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("X-Pylon-Token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no token header",
		})
		return
	}

	claims, err := a.JwtAuth.AuthorizeToken(tokenString)

	if err != nil {
		if err == ErrTokenInValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		} else if err == ErrTokenExpired {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "token is expired",
			})
			return
		} else if err == ErrTokenBlacklist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "token is expired",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}
	}

	userTokens, err := a.d.GetAllUserToken(claims.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	returnPayload := []*UserToken{}
	for _, token := range *userTokens {
		returnPayload = append(returnPayload, &UserToken{
			ExpireAt:   token.ExpireAt,
			Type:       token.Type,
			DeviceName: token.DeviceName,
		})
	}

	c.JSON(http.StatusOK, returnPayload)
}
