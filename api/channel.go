package api

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
)

type Database interface {
	CreateChannel(channel *model.Channel) error
	RemoveChannel(channelId uint) error
	GetChannelsByUserUUID(uuid string) (*[]model.Channel, error)
	InjectUserByChannelId(user *model.User, channelId uint) error
	RemoveUserByChannelId(user *model.User, channelId uint) error
}

type ChannelAPI struct {
	DB Database
}

type createChannelPayload struct {
	name    string
	members []string
}

type ChannelPayload struct {
	ChannelId uint
}

func (a *ChannelAPI) createChannelModel(*createChannelPayload) (*model.Channel, error) {
	return nil, nil
}

func (a *ChannelAPI) CreateChannel(c *gin.Context) {

}

func (a *ChannelAPI) RemoveChannel(c *gin.Context) {
	var channelPayload ChannelPayload
	err := c.BindJSON(&channelPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	err = a.DB.RemoveChannel(channelPayload.ChannelId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *ChannelAPI) GetChannelIds(c *gin.Context) {
	claims := c.MustGet("token").(auth.AuthTokenClaims)
	uuid := claims.UserUUID
	channels, err := a.DB.GetChannelsByUserUUID(uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	channelId := make([]uint, 0)
	for _, channel := range *channels {
		channelId = append(channelId, channel.Id)
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"channel": channelId,
	})
}

// authentication for join channel implement needed
func (a *ChannelAPI) JoinChannel(c *gin.Context) {
	// claims := c.MustGet("token").(auth.AuthTokenClaims)
	// uuid := claims.UserUUID
}
