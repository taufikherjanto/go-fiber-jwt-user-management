package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-fiber-user-management/model"
)

// instance database postgres
var DB *gorm.DB

// connect to postgres
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect DB")
	}

	fmt.Println("Success connect to DB")

	//Run migration DB
	err = DB.AutoMigrate(&model.User{}, &model.RevokedToken{})
	if err != nil {
		panic("Failed to run migration DB")
	}

	fmt.Println("Migration DB successfully")

}
