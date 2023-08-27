package database

import (
	"time"

	"github.com/PylonSchema/server/model"
	"github.com/PylonSchema/server/try"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Email    string
	Password string
}

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

func (d *GormDatabase) GetOriginUser(email string, password string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("email = ?", email).Find(user).Error
	if err != nil {
		return nil, err
	}
	origin := new(model.Origin)
	err = d.DB.Where("uuid = ?", user.UUID).Find(origin).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(origin.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, err
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
	err := d.DB.Where("email = ?", email).First(user).Error
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

func (d *GormDatabase) SetUserTokenPair(uuid uuid.UUID, expireAt time.Time, tokenString string) error {
	err := d.DB.Create(&model.UserTokenPair{
		UUID:     uuid,
		ExpireAt: expireAt,
		Token:    tokenString,
	}).Error
	return err
}

func (d *GormDatabase) GetAllUserToken(uuid uuid.UUID) (*[]model.UserTokenPair, error) {
	queryUserTokenPairs := []model.UserTokenPair{}
	err := d.DB.Where("uuid = ?", uuid).Find(&queryUserTokenPairs).Error
	if err != nil {
		return &queryUserTokenPairs, err
	}
	validUserTokenPairs := []model.UserTokenPair{}
	for _, userTokenPair := range queryUserTokenPairs {
		if userTokenPair.ExpireAt.Unix() < time.Now().Unix() {
			userToken := &UserToken{
				id: userTokenPair.ID,
				d:  d,
			}
			go try.TryN(userToken, 3)
			continue
		}
		validUserTokenPairs = append(validUserTokenPairs, userTokenPair)
	}
	return &validUserTokenPairs, nil
}

type UserToken struct {
	id uint
	d  *GormDatabase
}

func (r *UserToken) Run() error {
	return r.removeUserToken()
}

func (r *UserToken) Fail() error {
	return nil
}

func (r *UserToken) removeUserToken() error {
	userTokenPair := model.UserTokenPair{
		ID: r.id,
	}
	err := r.d.DB.Delete(&userTokenPair).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
