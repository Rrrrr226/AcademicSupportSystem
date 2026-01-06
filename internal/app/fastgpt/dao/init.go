package dao

import (
	"gorm.io/gorm"
)

var (
	Fastgpt = &fastgpt{}
)

func InitPG(db *gorm.DB) error {
	err := Fastgpt.Init(db)
	if err != nil {
		return err
	}

	return err
}
