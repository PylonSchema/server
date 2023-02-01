package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	Username    string `gorm:"primaryKey"`
	AccountType int
	UUID        uuid.UUID
	SecretUUID  uuid.UUID
	Email       string    `gorm:"unique;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Sse struct {
	SecretUUID uuid.UUID
	Salt       string
	Password   string
}

type Social struct {
	//secret_uuid   uuid.UUID
	SocialType   int
	Id           string
	AccessToken  string
	RefreshToken string
}

type Auth struct {
	Sse    Sse    `gorm:"embedded"`
	Social Social `gorm:"embedded"`
}

type Member struct {
	User User `gorm:"embedded"`
	Auth Auth `gorm:"embedded"`
}
