package api

import (
	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/gin-gonic/gin"
)

type Database interface {
	CreateChannel(channel *model.Channel) error
	GetChannelsByUserUUID(uuid string) ([]*model.Channel, error)
	InjectUserToChannelId(user *model.User, channelId int) error
	RemoveUserFromChannelId(user *model.User, channelId int) error
}

type ChannelAPI struct {
	DB Database
}

func (a *ChannelAPI) CreateChannel(c *gin.Context) {
	claims := c.MustGet("token").(auth.AuthTokenClaims)
	uuid := claims.UserUUID
	println(uuid)
}

func (a *ChannelAPI) GetChannelIds(c *gin.Context) {

}

func (a *ChannelAPI) JoinChannel(c *gin.Context) {

}

func (a *ChannelAPI) RemoveChannel(c *gin.Context) {

}
