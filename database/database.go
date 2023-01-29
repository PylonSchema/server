package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	fmt.Println("Connecting to Mysql Database")
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, err := d.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)

	// MySQL 서버, OS 또는 기타 미들웨어에 의해 연결이 종료되기 전에 안전하게 드라이버로 연결을 종료했는지 확인하는 데 필요
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	db = d
}

func GetDB() *gorm.DB {
	return db
}
