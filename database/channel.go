package database

import "github.com/PylonSchema/server/model"

type Channel struct {
}

func (d *GormDatabase) CreateChannel(channel *model.Channel) error {
	return d.DB.Create(channel).Error
}

func (d *GormDatabase) RemoveChannel(channelId uint) error {
	err := d.DB.Where(&model.Channel{}, channelId).Error
	return err
}

func (d *GormDatabase) GetChannelsByUserUUID(uuid string) (*[]model.Channel, error) {
	var channels []model.Channel
	err := d.DB.Where("uuid = ?", uuid).Find(&channels).Error
	return &channels, err
}

func (d *GormDatabase) InjectUserByChannelId(user *model.User, channelId uint) error {
	channelMember := model.ChannelMember{
		ChannelId:  uint(channelId),
		UUID:       user.UUID,
		SecretUUID: user.SecretUUID,
	}
	err := d.DB.Create(&channelMember).Error
	return err
}

func (d *GormDatabase) RemoveUserByChannelId(user *model.User, channelId uint) error {
	err := d.DB.Where("channel_id = ? AND secret_uuid = ?", channelId, user.SecretUUID).Delete(&model.ChannelMember{}).Error
	return err
}
