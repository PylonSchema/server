package origin

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/gin-gonic/gin"
)

func (a *AuthOriginAPI) LoginAccountHandler(c *gin.Context) {
	var loginPaylaod loginPaylaod
	err := c.BindJSON(&loginPaylaod)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "in-valid form error",
		})
		return
	}

	user, err := a.DB.GetOriginUser(loginPaylaod.Email, loginPaylaod.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	jp := auth.JwtPayload{
		UserUUID: user.UUID,
		Username: user.Username,
	}
	jwtTokenString, err := a.JwtAuth.GenerateJWT(&jp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"token":   jwtTokenString,
		"expire":  60 * 60 * 24 * 90,
	})
}
