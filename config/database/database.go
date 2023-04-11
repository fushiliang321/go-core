package database

import (
	"gorm.io/gorm"
	"time"
)

type Pool struct {
	MaxOpenConns int           //打开数据库连接的最大数量
	MaxIdleConns int           //连接池中空闲连接的最大数量
	MaxIdleTime  time.Duration //连接最大闲置时间
	MaxLifetime  time.Duration //连接最大可复用时间
	Heartbeat    time.Duration //心跳
}

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
	Pool     *Pool
	Settings *gorm.Config
}

var data = &Database{}

func Set(config *Database) {
	data = config
}

func Get() *Database {
	return data
}
