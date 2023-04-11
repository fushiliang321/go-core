package db

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/database"
	"github.com/fushiliang321/go-core/db/model"
	"github.com/fushiliang321/go-core/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var _db *gorm.DB

func db() *gorm.DB {
	if _db != nil {
		return _db
	}
	var err error
	config := database.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset)
	if config.Settings == nil {
		config.Settings = &gorm.Config{}
	}
	_db, err = gorm.Open(mysql.Open(dsn), config.Settings)
	if err != nil {
		exception.Listener("open db", err)
		return _db
	}
	if config.Pool == nil {
		return _db
	}
	pool := config.Pool
	mysqlDb, _ := _db.DB()
	if pool.MaxOpenConns > 0 {
		mysqlDb.SetMaxOpenConns(pool.MaxOpenConns)
	}
	if pool.MaxIdleConns > 0 {
		mysqlDb.SetMaxIdleConns(pool.MaxIdleConns)
	}
	if pool.MaxIdleTime > 0 {
		mysqlDb.SetConnMaxIdleTime(pool.MaxIdleTime)
	}
	if pool.MaxIdleTime > 0 {
		mysqlDb.SetConnMaxLifetime(pool.MaxIdleTime)
	}
	if pool.Heartbeat > 0 {
		go func() {
			for {
				time.Sleep(pool.Heartbeat)
				mysqlDb.Ping()
			}
		}()
	}
	return _db
}

func Model[T any]() *model.Model[T] {
	var t T
	return SetDb[T](db().Model(t))
}

func Db() *gorm.DB {
	return db()
}

func SetDb[T any](db *gorm.DB) *model.Model[T] {
	return &model.Model[T]{
		Db: db,
	}
}
