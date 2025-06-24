package model

import "gorm.io/gorm"

// 数据库模型

type Subject struct {
	gorm.Model

	SubjectName string `gorm:"type:varchar(20);not null"`

	SubjectLink string `gorm:"not null"`
}
