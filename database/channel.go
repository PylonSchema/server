package database

import (
	"github.com/PylonSchema/server/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *Database) CreateChannel(channel *model.Channel) error {
	return d.DB.Create(channel).Error
}

func (d *Database) RemoveChannel(channelId uint) error {
	err := d.DB.Where(&model.Channel{}, channelId).Error
	return err
}

func (d *Database) GetChannelsByUserUUID(uuid uuid.UUID) (*[]model.ChannelMember, error) {
	var channelMembers []model.ChannelMember
	err := d.DB.Where("uuid = ?", uuid).Find(&channelMembers).Error
	return &channelMembers, err
}

func (d *Database) InjectUserByChannelId(user *model.User, channelId uint) error {
	channelMember := model.ChannelMember{
		ChannelId: channelId,
		UUID:      user.UUID,
	}
	err := d.DB.Create(&channelMember).Error
	return err
}

func (d *Database) RemoveUserByChannelId(user *model.User, channelId uint) error {
	err := d.DB.Where("channel_id = ? AND uuid = ?", channelId, user.UUID).Delete(&model.ChannelMember{}).Error
	return err
}

func (d *Database) IsUserInChannelByUUID(userUUID uuid.UUID, channelId uint) (bool, error) {
	channelMember := new(model.ChannelMember)
	err := d.DB.Where("channel_id = ? AND uuid = ?", channelId, userUUID.String()).Find(channelMember).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if channelMember.UUID == userUUID {
		return true, nil
	}
	return false, nil
}
