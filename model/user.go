package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	ID          uint `gorm:"primaryKey;autoIncrement:true"`
	Username    string
	AccountType int
	UUID        uuid.UUID `gorm:"type:varchar(36);uniqueIndex;not null"`
	Email       string    `gorm:"size:64;uniqueIndex;not null;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Origin struct {
	UUID     uuid.UUID
	Password string
}

type Social struct {
	UUID         uuid.UUID
	SocialType   int
	SocialId     string
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	AccessToken  string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
}

type UserTokenPair struct {
	ID         uint      `gorm:"primaryKey;autoIncrement:true"`
	UUID       uuid.UUID `gorm:"index;not null"`
	ExpireAt   time.Time `gorm:"index"`
	Token      string    `gorm:"not null"`
	Type       int       `gorm:"not null"`
	DeviceName string    `gorm:"type:varchar(255);not null"`
}
