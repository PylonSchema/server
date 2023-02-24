package database

import (
	"github.com/PylonSchema/server/model"
	"gorm.io/gorm"
)

// CreateUser
func (d *GormDatabase) CreateUser(user *model.User) error {
	return d.DB.Create(user).Error
}

func (d *GormDatabase) CreateOriginUser(user *model.User, origin *model.Origin) error {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(origin).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateUser
func (d *GormDatabase) UpdateUser(user *model.User) error {
	return d.DB.Save(user).Error
}

func (d *GormDatabase) CreateSocial(social *model.Social) error {
	return d.DB.Create(social).Error
}

func (d *GormDatabase) IsEmailUsed(email string) (bool, error) {
	user := new(model.User)
	err := d.DB.Where("email = ?", email).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if user.Email == email {
		return true, nil
	}
	return false, nil
}

func (d *GormDatabase) GetUserFromSocialByEmail(email string, socialType int) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("email = ? AND account_type >= ?", email, socialType).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
