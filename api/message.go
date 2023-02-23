package api

import (
	"net/http"

	"github.com/PylonSchema/server/api/gateway"
	"github.com/PylonSchema/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessagePayload struct {
	ChannelId uint   `json:"channelid" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type MessageDatabase interface {
	IsUserInChannelByUUID(userUUID uuid.UUID, channelID uint) (bool, error)
}

type MessageAPI struct {
	g  *gateway.Gateway
	DB MessageDatabase
}

func (a *MessageAPI) CreateMessage(c *gin.Context) {
	messagePayload := &MessagePayload{}
	err := c.Bind(&messagePayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "bind json error",
		})
		return
	}
	claims := c.MustGet("claims").(auth.AuthTokenClaims)
	find, err := a.DB.IsUserInChannelByUUID(claims.UserUUID, messagePayload.ChannelId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "is user in channel by uuid error, db",
		})
		return
	}
	if !find {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "this user not in channel error",
		})
		return
	}
	data := map[string]interface{}{
		"data": messagePayload.Content,
	}
	message := &gateway.Message{
		Op: gateway.MessageData,
		D:  data,
	}
	err = a.g.Boardcast(messagePayload.ChannelId, message)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "boardcast message to gateway error",
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"message": "message successfully created",
	})
}
