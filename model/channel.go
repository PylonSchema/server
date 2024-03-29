package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type Channel struct {
	Id        uint             `gorm:"primaryKey;autoIncrement:true"`
	Name      string           `gorm:"type:varchar(225);not null"`
	UUID      uuid.UUID        `gorm:"type:varchar(36);not null"`
	Owner     uuid.UUID        `gorm:"type:varchar(36);not null"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
	Members   []*ChannelMember `gorm:"foreignKey:ChannelId;References:Id;constraint:OnDelete:CASCADE"`
}

type ChannelMember struct {
	ChannelId uint `gorm:"index"`
	UUID      uuid.UUID
}

type InvitationChannel struct {
	ChannelID       uint      `gorm:"index:not null"`
	InvitiationLink string    `gorm:"index;type:varchar(125);not null"` // random invitation uuid + random string (15)
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	ExpireType      int       `gorm:"index"`
}
