package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type Channel struct {
	Id        uint             `gorm:"primaryKey;autoIncrement:true"`
	Name      string           `gorm:"type:varchar(225);not null;default:null"`
	UUID      uuid.UUID        `gorm:"type:varchar(36);not null;default:null"`
	Owner     uuid.UUID        `gorm:"type:varchar(36);not null;default:null"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
	Members   []*ChannelMember `gorm:"foreignKey:ChannelId;References:Id;constraint:OnDelete:CASCADE"`
}

type ChannelMember struct {
	ChannelId uint `gorm:"index"`
	UUID      uuid.UUID
}

type InvitationChannel struct {
	ChannelUUID     uuid.UUID `gorm:"index:not null;default:null"`
	InvitiationUUID string    `gorm:"index;type:varchar(125);not null;default:null"` // random invitation uuid + random string (15)
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	ExpireAt        time.Time `gorm:"index"`
}
