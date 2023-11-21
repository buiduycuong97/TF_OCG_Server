package database_manager

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// const DB_USERNAME = "sql12662305"
// const DB_PASSWORD = "juB5hiHp6b"
// const DB_NAME = "sql12662305"
// const DB_HOST = "sql12.freesqldatabase.com"

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}

	// Đọc biến môi trường từ tệp .env
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Tạo chuỗi kết nối DSN từ biến môi trường
	dsn := dbUsername + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	return db
}
