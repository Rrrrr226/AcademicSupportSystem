package model

import (
	"HelpStudent/internal/model"

	"gorm.io/gorm"
)

// 数据库模型

type Managers struct {
	model.Base
	StaffId   string         `gorm:"uniqueIndex:idx_staff_id;size:19"`
	DeletedAt gorm.DeletedAt `gorm:"uniqueIndex:idx_staff_id" json:"-"`
}
