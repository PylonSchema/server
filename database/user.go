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

//email 중복 확인
func (d *GormDatabase) IsEmailUsed(email string) (*model.User) {
	user := new(model.User)
	result := d.DB.First(&user, "email = ?", email)
	if(result.email == email)	{
		return true
	}
	else return false
}

func 