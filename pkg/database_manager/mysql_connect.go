package database_manager

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// const DB_USERNAME = "sql12662305"
// const DB_PASSWORD = "juB5hiHp6b"
// const DB_NAME = "sql12662305"
// const DB_HOST = "sql12.freesqldatabase.com"
const DB_USERNAME = "root"
const DB_PASSWORD = "cuong123"
const DB_NAME = "e-commerce"
const DB_HOST = "localhost"
const DB_PORT = "3306"

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Error connecting to database : error=%v", err)
	}

	return db
}
