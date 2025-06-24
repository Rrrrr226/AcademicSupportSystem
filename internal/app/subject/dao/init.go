package dao

import (
	"gorm.io/gorm"
)

var (
	Subject = &subject{}
)

func InitPG(db *gorm.DB) error {
	err := Subject.Init(db)
	if err != nil {
		return err
	}

	return err
}
