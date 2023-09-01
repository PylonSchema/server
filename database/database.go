package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PylonSchema/server/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	SQLConfig   *SQLConfig
	NOSQLConfig *NOSQLConfig
}

type SQLConfig struct {
	Username string
	Password string
	Address  string
	Port     string
}

type NOSQLConfig struct {
	Hosts []string
	Port  string
}

type Database struct {
	DB       *gorm.DB
	ScyllaDB *gocqlx.Session
}

var db *Database

func New(databaseConfig *DatabaseConfig) (*Database, error) {
	if db != nil {
		return db, nil
	}

	// sql database connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/pylon?charset=utf8mb4&parseTime=True&loc=Local",
		databaseConfig.SQLConfig.Username,
		databaseConfig.SQLConfig.Password,
		databaseConfig.SQLConfig.Address,
		databaseConfig.SQLConfig.Port)
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

	// scylla database connect
	cluster := gocql.NewCluster(databaseConfig.NOSQLConfig.Hosts...)
	if databaseConfig.NOSQLConfig.Port != "" {
		cluster.Port, err = strconv.Atoi(databaseConfig.NOSQLConfig.Port)
		if err != nil {
			return nil, err
		}
	}
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}

	db = &Database{DB: d, ScyllaDB: &session}

	return db, nil
}

func (g *Database) AutoMigration() error {

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
