package dao

import (
	"HelpStudent/internal/app/managers/model"
	"gorm.io/gorm"
)

type managers struct {
	*gorm.DB
}

func (u *managers) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Managers{})
}
