package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type Channel struct {
	Id        uint `gorm:"primaryKey;autoIncrement:true"`
	Name      string
	UUID      uuid.UUID
	Owner     uuid.UUID
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
	Members   []*ChannelMember `gorm:"foreignKey:ChannelId;References:Id;constraint:OnDelete:CASCADE"`
}

type ChannelMember struct {
	ChannelId uint `gorm:"index"`
	UUID      uuid.UUID
}

type InvitationChannel struct {
	ChannelUUID     uuid.UUID `gorm:"index"`
	InvitiationUUID string    `gorm:"index;type:varchar(125)"` // random invitation uuid + random string (15)
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	ExpireAt        time.Time
}
