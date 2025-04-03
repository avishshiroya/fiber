package config

import (
	"fiber/app/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=2034 dbname=fiber port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	fmt.Println("DB Connected successfully")

	// Run auto-migration
	err = db.AutoMigrate(&models.User{}, &models.Auth{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	DB = db
}
