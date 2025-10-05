package config

import (
	"fmt"
	"log"
	"os"

	"sanbercode-golang-batch-70-final-project/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// migrate otomatis
	db.AutoMigrate(&models.Role{}, &models.User{}, &models.LetterType{}, &models.Letter{}, &models.Setting{})

	DB = db
}
