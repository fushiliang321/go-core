package db

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/database"
	"github.com/fushiliang321/go-core/db/model"
	"github.com/fushiliang321/go-core/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	mysqlDb, _ := _db.DB()
	mysqlDb.SetMaxOpenConns(10)
	mysqlDb.SetMaxIdleConns(1)
	return _db
}

func Model[T any]() *model.Model[T] {
	var t T
	_db1 := db().Model(t)
	return &model.Model[T]{
		Db: _db1,
	}
}

func Db() *gorm.DB {
	return db()
}
