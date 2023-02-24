package origin

// this package manage origin platform account, Pylon

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Database interface {
	CreateOriginUser(user *model.User, origin *model.Origin) error
	IsEmailUsed(email string) (bool, error)
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

	isValid := createPayload.isValid()
	if !isValid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "user create form error, in-valid request",
		})
		return
	}

	hashedPassword, err := createPayload.hashing()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "hashing password error",
		})
		return
	}
	userModel, originModel, err := createPayload.createModel(hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "create model error",
		})
		return
	}

	isEmailUsed, err := a.DB.IsEmailUsed(createPayload.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "query db error",
		})
		return
	}
	if isEmailUsed {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "email already used",
		})
		return
	}

	err = a.DB.CreateOriginUser(userModel, originModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "create user in db error (transaction error)",
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (a *AuthOriginAPI) LoginAccountHandler(c *gin.Context) {

}

func (c *createPayload) createModel(hashedPassword string) (*model.User, *model.Origin, error) {
	userUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}
	return &model.User{
			Username:    c.Username,
			AccountType: 1,
			UUID:        userUUID,
			Email:       c.Email,
		}, &model.Origin{
			UUID:     userUUID,
			Password: hashedPassword,
		}, nil
}

func (c *createPayload) hashing() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (c *createPayload) isValid() bool {
	if len(c.Id) < 6 {
		return false
	}
	if len(c.Password) < 10 {
		return false
	}
	if !vaildEmail(c.Email) {
		return false
	}
	return true
}
