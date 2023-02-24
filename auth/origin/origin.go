package origin

// this package manage origin platform account, Pylon

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (a *AuthOriginAPI) createModel(username string, email string, password string) (*model.User, *model.Origin, error) {
	userUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}
	return &model.User{
			Username:    username,
			AccountType: 1,
			UUID:        userUUID,
			Email:       email,
		}, &model.Origin{
			UUID:     userUUID,
			Password: password,
		}, nil
}
