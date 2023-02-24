package origin

// this package manage origin platform account, Pylon

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/gin-gonic/gin"
)

type Database interface {
}

type AuthOriginAPI struct {
	DB      Database
	JwtAuth *auth.JwtAuth
}

type createPayload struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func New(db Database, jwtAuth *auth.JwtAuth) *AuthOriginAPI {
	return &AuthOriginAPI{
		DB:      db,
		JwtAuth: jwtAuth,
	}
}

func (a *AuthOriginAPI) CreateAccountHandler(c *gin.Context) {
	var createPayload createPayload
	err := c.BindJSON(&createPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "bind json error",
		})
		return
	}

}

func (a *AuthOriginAPI) createAccount() {

}

func (a *AuthOriginAPI) createTransaction() {
}
