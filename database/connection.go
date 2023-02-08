package database

import (
	"github.com/Akmdvtm/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Подключаемся к БД с помощью gorm.Open(1: Sql driver(dsn: "root:pass@/path"), 2: Config) db, err
	connection, err := gorm.Open(mysql.Open("root:12345Acdc@/go_auth"), &gorm.Config{})
	if err != nil {
		panic("could not connect to db")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
