package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	ID          uint `gorm:"primaryKey;autoIncrement:true"`
	Username    string
	AccountType int
	UUID        uuid.UUID `gorm:"index:"`
	Email       string    `gorm:"unique;not null"`
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
	ID       uint      `gorm:"primaryKey;autoIncrement:true"`
	UUID     uuid.UUID `gorm:"index:"`
	ExpireAt time.Time `gorm:"index:"`
	Token    string
	Type     int
}
