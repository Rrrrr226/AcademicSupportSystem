package dao

import (
	"HelpStudent/internal/app/users/model"
	"gorm.io/gorm"
)

type users struct {
	*gorm.DB
}

func (u *users) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Users{})
}
