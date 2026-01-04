package model

import (
	"HelpStudent/internal/model"
	"time"
)

// UserSubject 用户-科目关联表
type UserSubject struct {
	model.Base
	UserId      string    `gorm:"type:char(26);not null;index:idx_user_subject,unique" json:"-"`
	StaffId     string    `gorm:"type:varchar(19);not null;index" json:"staff_id"`
	SubjectName string    `gorm:"type:varchar(50);not null;index:idx_user_subject,unique" json:"subject_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UserSubject) TableName() string {
	return "user_subjects"
}
