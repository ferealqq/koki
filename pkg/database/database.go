package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var (
	KeyTable = "key_events"
)

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
// test connection
func TConn(dbFile string) *gorm.DB {
	if db == nil {
		d, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = d 
	}
	return db
}
//SELECT COUNT(Id),strftime ('%H',timestamp) hour
//FROM T
//GROUP BY strftime ('%H',timestamp)