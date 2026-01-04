package model

import (
	"time"

	"gorm.io/gorm"
)

// 数据库模型

type Subject struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	SubjectName string         `gorm:"type:varchar(20);not null" json:"subject_name"`
	SubjectLink string         `gorm:"not null" json:"subject_link"`
}
