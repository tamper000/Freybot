package database

import (
	"github.com/tamper000/freybot/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LoadDatabase(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Message{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
