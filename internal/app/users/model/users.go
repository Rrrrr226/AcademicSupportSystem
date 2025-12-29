package model

import (
	"HelpStudent/pkg/utils"
	"gorm.io/gorm"
)

type Users struct {
	Id       string `gorm:"primary_key;type:char(26)"`
	StaffId  string `gorm:"uniqueIndex;size:19"`
	Name     string
	Username string `gorm:"size:50"`
	Password string
	Email    string
	Phone    string

	NeedSubjects   []string `gorm:"-"`                              // 不直接映射到数据库字段
	NeedSubjectsDB string   `gorm:"column:need_subjects;type:text"` // 实际数据库存储字段
}

func (u *Users) BeforeCreate(*gorm.DB) error {
	u.Id = utils.GenUUID()
	return nil
}
