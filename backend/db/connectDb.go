package db

import (
	"backend/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func ConnectDb() {
	Db, err = gorm.Open(mysql.Open(os.Getenv("DB_URI")), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
}

func MigrateDb() {
	err := Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error aauto migrating schema!", err)
	}
}
