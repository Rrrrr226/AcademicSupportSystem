package dao

import (
	"gorm.io/gorm"
)

var (
	Permission = &permission{}
)

func InitPG(db *gorm.DB) error {
	err := Permission.Init(db)
	if err != nil {
		return err
	}

	return err
}
