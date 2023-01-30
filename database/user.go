package database

import (
	"time"

	uuid "github.com/google/uuid"

	"path/filepath"
	"gorm.io/gorm"
	"C:\Users\Administrator\Desktop\DevStudy\GoLangChat\sse-chat-main\model"
)

// CreateUser
func (d *GormDatabase) CreateUser(user *model.User) error {
	return d.DB.Create(user).Error
}

// UpdateUser
func (d *GormDatabase) UpdateUser(user *model.User) error {
	return d.DB.Save(user).Error
}

// Query
func (d *GormDatabase) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("name = ?", name).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if user.userame == name {
		return user, err
	}
	return nil, err
}

func (d *GormDatabase) GetAuthByUUID(uuid uuid.UUID) (*model.Auth, error) {
	auth := new(model.Auth)
	err := d.DB.Where("secret_uuid = ?", uuid).Find(auth).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if auth.secret_uuid == uuid {
		return auth, err
	}
	return nil, err
}

//email 중복 확인
func (d *GormDatabase) IsEmailUsed(email string) (*model.User) {
	user := new(model.User)
	result := d.DB.First(&user, "email = ?", email)
	if(result.email == email)	{
		return truess
	}
	return false
}
