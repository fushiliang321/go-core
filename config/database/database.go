package database

import "gorm.io/gorm"

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
	Settings *gorm.Config
}

var data = &Database{}

func Set(config *Database) {
	data = config
}

func Get() *Database {
	return data
}
