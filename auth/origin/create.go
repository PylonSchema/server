package origin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type createPayload struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (c *createPayload) isValid() bool {
	if len(c.Password) < 10 {
		return false
	}
	if !vaildEmail(c.Email) {
		return false
	}
	return true
}

func (c *createPayload) hashing() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (a *AuthOriginAPI) CreateAccountHandler(c *gin.Context) {
	var createPayload createPayload
	err := c.BindJSON(&createPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
