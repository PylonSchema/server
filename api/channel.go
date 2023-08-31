package api

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChannelDatabase interface {
	CreateChannel(channel *model.Channel) error
	RemoveChannel(channelId uint) error
	GetChannelsByUserUUID(uuid uuid.UUID) (*[]model.ChannelMember, error)
	InjectUserByChannelId(user *model.User, channelId uint) error
	RemoveUserByChannelId(user *model.User, channelId uint) error
}

type ChannelGateway interface {
}

type ChannelAPI struct {
	d ChannelDatabase
	g ChannelGateway
}

type createChannelPayload struct {
	Name string `json:"name" binding:"required"`
}

type ChannelPayload struct {
	ChannelId uint
}

func NewChannelAPI(database ChannelDatabase, gateway ChannelGateway) *ChannelAPI {
	return &ChannelAPI{
		d: database,
		g: gateway,
	}
}

func (a *ChannelAPI) CreateChannelHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.AuthTokenClaims)

	var payload createChannelPayload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "bind json error"})
		return
	}
	channelUUID, err := uuid.NewRandom()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "create channel model error"})
		return
	}
	channelModel := &model.Channel{
		Name:  payload.Name,
		UUID:  channelUUID,
		Owner: claims.UserUUID,
	}
	err = a.d.CreateChannel(channelModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "create channel error"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *ChannelAPI) RemoveChannelHandler(c *gin.Context) {
	var channelPayload ChannelPayload
	err := c.BindJSON(&channelPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	err = a.d.RemoveChannel(channelPayload.ChannelId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *ChannelAPI) GetChannelIdsHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.AuthTokenClaims)
	uuid := claims.UserUUID
	channels, err := a.d.GetChannelsByUserUUID(uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	channelId := make([]uint, 0)
	for _, channel := range *channels {
		channelId = append(channelId, channel.ChannelId)
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"channel": channelId,
	})
}

// authentication for join channel implement needed
func (a *ChannelAPI) JoinChannelHandler(c *gin.Context) {
	// claims := c.MustGet("token").(auth.AuthTokenClaims)
	// uuid := claims.UserUUID
}

// remove user from channel
func (a *ChannelAPI) RemoveUserHandler(c *gin.Context) {
	//
}
