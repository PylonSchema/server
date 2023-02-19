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
	GetChannelsByUserUUID(uuid uuid.UUID) (*[]model.Channel, error)
	InjectUserByChannelId(user *model.User, channelId uint) error
	RemoveUserByChannelId(user *model.User, channelId uint) error
}

type ChannelAPI struct {
	DB ChannelDatabase
}

type createChannelPayload struct {
	name    string
	members []string
}

type ChannelPayload struct {
	ChannelId uint
}

func (a *ChannelAPI) createChannelModel(payload *createChannelPayload, owner uuid.UUID) (*model.Channel, error) {
	channelUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	membersSet := make(map[uuid.UUID]struct{})
	members := make([]model.ChannelMember, 0)
	for _, member := range payload.members {
		memberUUID, err := uuid.Parse(member)
		if err != nil {
			continue
		}
		if memberUUID == owner {
			continue
		}
		_, found := membersSet[memberUUID]
		if !found {
			membersSet[memberUUID] = struct{}{}
			members = append(members, model.ChannelMember{
				UUID: memberUUID,
			})
		}
	}
	members = append(members, model.ChannelMember{
		UUID: owner,
	})
	model := &model.Channel{
		Name:    payload.name,
		UUID:    channelUUID,
		Owner:   owner,
		Members: members,
	}
	return model, nil
}

func (a *ChannelAPI) CreateChannel(c *gin.Context) {
	claims := c.MustGet("token").(auth.AuthTokenClaims)
	uuid := claims.UserUUID

	var payload createChannelPayload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "bind json error"})
		return
	}

	payload.members = append(payload.members, uuid.String())
	channelModel, err := a.createChannelModel(&payload, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "create channel model error"})
		return
	}
	err = a.DB.CreateChannel(channelModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "create channel error"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "ok"})
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

// remove user from channel
func (a *ChannelAPI) RemoveUser(c *gin.Context) {
	//
}
