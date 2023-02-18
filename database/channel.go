package database

import "github.com/PylonSchema/server/model"

type Channel struct {
}

func (d *GormDatabase) CreateChannel(channel *model.Channel) error {
	return nil
}

func (d *GormDatabase) GetChannelsByUserUUID(uuid string) ([]*model.Channel, error) {
	return nil, nil
}

func (d *GormDatabase) InjectUserToChannelId(user *model.User, channelId int) error {
	return nil
}

func (d *GormDatabase) RemoveUserFromChannelId(user *model.User, channelId int) error {
	return nil
}
