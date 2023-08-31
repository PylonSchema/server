package database

import (
	"fmt"
	"time"

	"github.com/PylonSchema/server/model"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDatabase struct {
	DB *gorm.DB
}

var db *GormDatabase

func New(username string, password string, address string, port string) (*GormDatabase, error) {
	if db != nil {
		return db, nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/pylon?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, port)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqldb, err := d.DB()

	if err != nil {
		return nil, err
	}

	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(0)
	sqldb.SetConnMaxLifetime(time.Minute * 3)

	db = &GormDatabase{DB: d}

	return db, nil
}

func (g *GormDatabase) AutoMigration() error {

	err := g.DB.AutoMigrate(
		&model.User{}, &model.Origin{}, &model.Social{},
		&model.Channel{}, &model.ChannelMember{}, &model.InvitationChannel{},
		&model.RefreshToken{}, &model.UserTokenPair{},
	)
	if err != nil {
		panic(err)
	}

	return nil
}
