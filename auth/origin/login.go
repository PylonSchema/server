package origin

import (
	"net/http"
	"strings"
	"time"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
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

	userAgentHeader := c.GetHeader("User-Agent")
	if userAgentHeader == "" {
		userAgentHeader = "Unknown"
	}

	userAgentType := strings.Split(userAgentHeader, "/")[0]

	expireSec := 60 * 30
	deviceType := 0
	if userAgentType == "PylonMobile" {
		expireSec = 60 * 60 * 2 // 2 hour
		deviceType = 1
	} else if userAgentType == "PylonDesktop" {
		expireSec = 60 * 60 * 6 // 6 hour
		deviceType = 2
	}

	expireAt := time.Now().Add(time.Second * time.Duration(expireSec))

	userTokenPair := &model.UserTokenPair{
		UUID:       user.UUID,
		ExpireAt:   expireAt,
		Token:      jwtTokenString,
		Type:       deviceType,
		DeviceName: userAgentHeader,
	}

	err = a.DB.SetUserTokenPair(userTokenPair)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, auth.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"token":   jwtTokenString,
		"expire":  expireSec,
	})
}
