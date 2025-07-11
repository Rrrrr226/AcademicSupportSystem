package dao

import (
	"gorm.io/gorm"
)

type permission struct {
	*gorm.DB
}

func (u *permission) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate()
}
