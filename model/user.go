package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	Username    string
	AccountType int
	UUID        uuid.UUID
	SecretUUID  uuid.UUID
	Email       string    `gorm:"unique;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Origin struct {
	SecretUUID uuid.UUID
	Salt       string
	Password   string
}

type Social struct {
	SecretUUID   uuid.UUID
	SocialType   int
	SocialId     string
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	UUID         uuid.UUID `gorm:"not null"`
	AccessToken  string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
}
