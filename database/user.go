package database

import (
	"github.com/devhoodit/sse-chat/model"
	"gorm.io/gorm"
)

// CreateUser
func (d *GormDatabase) CreateUser(user *model.User) error {
	return d.DB.Create(user).Error
}

// UpdateUser
func (d *GormDatabase) UpdateUser(user *model.User) error {
	return d.DB.Save(user).Error
}

func (d *GormDatabase) CreateSocial(social *model.Social) error {
	return d.DB.Create(social).Error
}

// Query
func (d *GormDatabase) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("name = ?", name).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if user.Username == name {
		return user, err
	}
	return nil, err
}

// func (d *GormDatabase) GetAuthByUUID(uuid uuid.UUID) (*model.Auth, error) {
// 	auth := new(model.Auth)
// 	err := d.DB.Where("secret_uuid = ?", uuid).Find(auth).Error
// 	if err == gorm.ErrRecordNotFound {
// 		err = nil
// 	}
// 	if auth.secret_uuid == uuid {
// 		return auth, err
// 	}
// 	return nil, err
// }

// email 중복 확인
func (d *GormDatabase) IsEmailUsed(email string) bool {
	user := new(model.User)
	err := d.DB.Where("email = ?", email).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	if user.Email == email {
		return true
	}
	return false
}
