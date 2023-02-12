package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	Id        uint `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Members   []ChannelMember
}

type ChannelMember struct {
	gorm.Model
	Id         int
	UUID       uuid.UUID
	SecretUUID uuid.UUID
}
