package dao

import (
	"HelpStudent/internal/app/fastgpt/model"

	"gorm.io/gorm"
)

type fastgpt struct {
	*gorm.DB
}

func (u *fastgpt) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Fastgpt{})
}
