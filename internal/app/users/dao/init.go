package dao

import (
	"errors"

	"gorm.io/gorm"
)

var (
	Users = &users{DB: nil}
)

func InitPG(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	return Users.Init(db)
}
