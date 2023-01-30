package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	username     string `gorm:"primaryKey"`
	account_type int
	uuid         uuid.UUID
	secret_uuid  uuid.UUID
	email        string    `gorm:"unique;not null"`
	created_at   time.Time `gorm:"autoCreateTime"`
	updated_at   time.Time `gorm:"autoUpdateTime"`
}

type Sse struct {
	secret_uuid uuid.UUID
	salt        string
	pathword    string
}

type Social struct {
	//secret_uuid   uuid.UUID
	social_type   int
	id            string
	access_token  string
	refresh_token string
}

type Auth struct {
	Sse    Sse    `gorm:"embedded"`
	Social Social `gorm:"embedded"`
}

type Member struct {
	User User `gorm:"embedded"`
	Auth Auth `gorm:"embedded"`
}
