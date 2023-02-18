package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type Channel struct {
	Id        uint `gorm:"primaryKey;autoIncrement:true"`
	Name      string
	UUID      uuid.UUID
	CreatedAt time.Time       `grom:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	Members   []ChannelMember `gorm:"foreignKey:ChannelId;References:Id;OnUpdate:CASCADE"`
}

type ChannelMember struct {
	ChannelId  uint `gorm:"index"`
	UUID       uuid.UUID
	SecretUUID uuid.UUID
}
