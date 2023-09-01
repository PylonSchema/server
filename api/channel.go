package api

import (
	"net/http"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/PylonSchema/server/pylontype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChannelDatabase interface {
	CreateChannel(channel *model.Channel) error
	RemoveChannel(channelId uint) error
	GetChannelsByUserUUID(uuid uuid.UUID) (*[]model.ChannelMember, error)
	InjectUserByChannelId(user *model.User, channelId uint) error
	RemoveUserByChannelId(user *model.User, channelId uint) error
	GetUserRoleInChannelByUUID(userUUID uuid.UUID, channelId uint) (int, error)
	CreateChannelInvitationLink(channel_id uint, link_type int) (string, error)
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
	channelMember := []*model.ChannelMember{{UUID: claims.UserUUID}}
	channelModel := &model.Channel{
		Name:    payload.Name,
		UUID:    channelUUID,
		Owner:   claims.UserUUID,
		Members: channelMember,
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

func (a *ChannelAPI) GetChannelInvitationLinkHandler(c *gin.Context) {
	//
}

type createChannelInvitationLink struct {
	ChannelUUID uint `json:"channel_uuid"`
	ExpireType  int  `json:"expire_type"` // dispoable, 1 hour, 1 day, 1 week, 1 month, permanent
}

func (a *ChannelAPI) CreateChannelInvitationLinkHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.AuthTokenClaims)
	var form createChannelInvitationLink
	err := c.BindJSON(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	if form.ExpireType < 0 || form.ExpireType > 6 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	userRole, err := a.d.GetUserRoleInChannelByUUID(claims.UserUUID, form.ChannelUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no permission to create link",
		})
		return
	}
	if userRole != pylontype.UserRoleOwner {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no permission to create link",
		})
		return
	}
	link, err := a.d.CreateChannelInvitationLink(form.ChannelUUID, form.ExpireType)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no permission to create link",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"link": link,
	})
}

func (a *ChannelAPI) RemoveChannelInvitationLinkHandler(c *gin.Context) {
	//
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
