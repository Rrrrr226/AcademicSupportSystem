package model

import "gorm.io/gorm"

// 数据库模型

type Users struct {
	gorm.Model

	StaffId string `gorm:"uniqueIndex;size:19"`
	Name    string

	NeedSubjects   []string `gorm:"-"`                              // 不直接映射到数据库字段
	NeedSubjectsDB string   `gorm:"column:need_subjects;type:text"` // 实际数据库存储字段
}
