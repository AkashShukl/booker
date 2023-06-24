package model

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("bookings.db"), &gorm.Config{})
	err = db.AutoMigrate(&Booking{})
	if err != nil {
		return err
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DEBUG: Init DB Completed")
	return nil
}

func GetDB() *gorm.DB {
	if db != nil {
		return db
	} else {
		InitDB()
		return db
	}
}
