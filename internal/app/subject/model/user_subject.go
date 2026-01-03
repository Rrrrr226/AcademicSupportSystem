package model

import (
	"HelpStudent/internal/model"
)

// UserSubject 用户-科目关联表
type UserSubject struct {
	model.Base
	UserId      string `gorm:"type:char(26);not null;index:idx_user_subject,unique"`
	StaffId     string `gorm:"type:varchar(19);not null;index"`
	SubjectName string `gorm:"type:varchar(50);not null;index:idx_user_subject,unique"`
}

// TableName 指定表名
func (UserSubject) TableName() string {
	return "user_subjects"
}
