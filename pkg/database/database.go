package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Conn() *gorm.DB{
	if db == nil {
		d, err := gorm.Open(sqlite.Open("logger/log.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = d 
	}
	return db
}

func Close() {
	db, err := db.DB()

	if err != nil {
		panic(err)
	}

	db.Close()
}
