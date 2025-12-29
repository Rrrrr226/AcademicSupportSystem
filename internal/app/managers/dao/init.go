package dao

import (
	"gorm.io/gorm"
)

var (
	Managers = &managers{}
)

func InitPG(db *gorm.DB) error {
	err := Managers.Init(db)
	if err != nil {
		return err
	}

	return err
}
