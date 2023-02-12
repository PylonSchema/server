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

func New(username string, password string, address string, port string) (*GormDatabase, error) {
	fmt.Println("Connecting to Mysql Database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/pylon?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, port)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := d.DB()

	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)

	// MySQL 서버, OS 또는 기타 미들웨어에 의해 연결이 종료되기 전에 안전하게 드라이버로 연결을 종료했는지 확인하는 데 필요
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	return &GormDatabase{DB: d}, nil
}

func (g *GormDatabase) AutoMigration() error {

	err := g.DB.AutoMigrate(&model.User{}, &model.Origin{}, &model.Social{})
	if err != nil {
		panic(err)
	}

	return nil
}
